package pm

import (
	"api/ipm"
	"fdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"iclient"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"scode"
	"status"
	"strconv"
	"strings"
	"sync"
	"time"
	"util"
)

const (
	PMDB             = "/var/lib/fal/pm.db"
	BucketPlayerList = "players"
	PlayerPrefix     = "player-"
	BucketKeyPid     = "pid"
	BucketKeyId      = "id"
	BucketKeyName    = "name"
	BucketKeyPwd     = "password"
	BucketKeyPort    = "port"
)

var (
	pm  *PlayerManager
	log = logrus.WithFields(logrus.Fields{"pkg": "pm"})
)

type PlayerManager struct {
	sync.RWMutex
	players map[string]*ipm.PlayerInfo
	DB      *fdb.FalDB
	isInit  bool
}

func NewPlayerManager() *PlayerManager {
	db := fdb.NewDB(PMDB)
	err := db.CreateBucket(BucketPlayerList)
	if err != nil {
		panic(err)
	}
	pm = &PlayerManager{
		players: make(map[string]*ipm.PlayerInfo),
		DB:      db,
	}
	return pm
}

// 初始化
func (m *PlayerManager) InitUpdate() error {
	m.Lock()
	defer m.Unlock()

	if m.isInit {
		panic("InitUpdate PM twice")
	}
	m.isInit = true

	if err := m.playerInit(); err != nil {
		return err
	}

	return m.scanPlayerInfo()
}

func (m *PlayerManager) scanPlayerInfo() error {
	match := []string{"fal", "player", "start"}
	pids, err := util.FindPids(match)
	if err != nil {
		return err
	}

	for _, pid := range pids {
		file := fmt.Sprintf("/proc/%d/cmdline", pid)
		output, err := ioutil.ReadFile(file)
		if err != nil {
			log.Errorf("read player info error:%s", err.Error())
			continue
		}
		cmdLine := strings.Split(string(output[:len(output)-1]), string(byte(0)))
		var id string
		var addr string
		for _, str := range cmdLine {
			if strings.HasPrefix(str, "--name=") {
				id = strings.TrimPrefix(str, "--name=")
			}
			if strings.HasPrefix(str, "--addr=") {
				addr = strings.TrimPrefix(str, "--addr=")
			}
		}

		if _, exist := m.players[id]; exist {
			_, portstr, _ := net.SplitHostPort(addr)
			m.players[id].Pid = int64(pid)
			m.players[id].Port, _ = strconv.ParseInt(portstr, 10, 64)
			if err := m.updatePlayerInfo(m.players[id]); err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

// 从数据库读入游戏数据
func (m *PlayerManager) playerInit() error {
	buckets, err := m.DB.GetAllBucket()
	if err != nil {
		desc := fmt.Sprintf("GetAllBucket error:%s", err)
		return status.NewStatusDesc(scode.PMDBOperateError, desc)
	}
	log.Debug(buckets)
	for _, b := range buckets {
		if strings.HasPrefix(b, PlayerPrefix) {

			kvs, err := m.DB.GetAllKV(b)
			if err != nil {
				desc := fmt.Sprintf("GetAllKV error:%s", err)
				return status.NewStatusDesc(scode.PMDBOperateError, desc)
			}

			port, _ := strconv.ParseInt(kvs[BucketKeyPort], 10, 64)
			pid, _ := strconv.ParseInt(kvs[BucketKeyPid], 10, 64)
			info := &ipm.PlayerInfo{
				Id:       kvs[BucketKeyId],
				Port:     port,
				Pid:      pid,
				Name:     kvs[BucketKeyName],
				Password: kvs[BucketKeyPwd],
			}
			m.players[info.Id] = info
		}
	}
	return nil
}

// 创建玩家
func (m *PlayerManager) CreatePlayer(req *ipm.PlayerCreateRequest) (*ipm.PlayerInfo, error) {
	m.Lock()
	defer m.Unlock()
	id := req.Pid
	if req.Pid != "" {
		_, exist := m.players[req.Pid]
		if exist {
			desc := fmt.Sprintf("player %s already exist", req.Pid)
			return nil, status.NewStatusDesc(scode.PlayerAlreadyExist, desc)
		}
	} else {
		id = uuid.New().String()
	}

	id = PlayerPrefix + id

	p := &ipm.PlayerInfo{
		Id:       id,
		Name:     req.Name,
		Password: req.Password,
	}
	kvs := make(map[string]string)
	kvs[BucketKeyId] = id
	kvs[BucketKeyName] = p.Name
	kvs[BucketKeyPwd] = p.Password
	data := map[string]map[string]string{id: kvs}
	m.players[id] = p

	m.Unlock()
	defer m.Lock()
	if err := m.DB.Put(id, req.Name, BucketPlayerList); err != nil {
		desc := fmt.Sprintf("Put key:%s value:%s bucket:%s error:%s", id, req.Name, BucketPlayerList, err)
		return nil, status.NewStatusDesc(scode.PMDBOperateError, desc)
	}
	if err := m.DB.PutBatch(data); err != nil {
		desc := fmt.Sprintf("PutBatch:%v error:%s", data, err)
		return nil, status.NewStatusDesc(scode.PMDBOperateError, desc)
	}
	return p, nil
}

// 分配端口
func (m *PlayerManager) assignPort(port int64) int64 {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		lis, err = net.Listen("tcp", ":")
	}
	if err != nil {
		panic(err)
	}
	lis.Close()
	return int64(lis.Addr().(*net.TCPAddr).Port)
}

// 删除玩家
func (m *PlayerManager) DeletePlayer(id string) error {
	m.Lock()
	defer m.Unlock()
	_, exist := m.players[id]
	if !exist {
		desc := fmt.Sprintf("player %s is not exist", id)
		return status.NewStatusDesc(scode.PlayerNotExist, desc)
	}
	delete(m.players, id)

	m.Unlock()
	defer m.Lock()
	if err := m.DB.Delete(id, BucketPlayerList); err != nil {
		log.Errorf("Delete key:%s bucket:%s error:%s", id, BucketPlayerList, err)
	}
	if err := m.DB.DeleteBucket(id); err != nil {
		log.Errorf("Delete bucket:%s error:%s", id, err)
	}
	return nil
}

// 获取玩家信息
func (m *PlayerManager) GetAllPlayerInfo() map[string]*ipm.PlayerInfo {
	m.Lock()
	defer m.Unlock()
	players := make(map[string]*ipm.PlayerInfo)
	for k, g := range m.players {
		players[k] = g
	}
	return players
}

// 启动进程
func (m *PlayerManager) StartPlayer(req *ipm.PlayerSignInRequest) (*ipm.PlayerInfo, error) {
	m.Lock()
	defer m.Unlock()
	p, exist := m.players[req.Pid]
	if !exist {
		desc := fmt.Sprintf("player %s is not exist", req.Pid)
		return nil, status.NewStatusDesc(scode.PlayerNotExist, desc)
	}
	if p.Name != req.Name || p.Password != req.Password {
		desc := fmt.Sprintf("player %s is name %s or password %s error", req.Pid, req.Name, req.Password)
		return nil, status.NewStatusDesc(scode.NameOrPasswordError, desc)
	}
	// 分配端口
	p.Port = m.assignPort(p.Port)

	args := []string{
		filepath.Base(os.Args[0]),
		"player",
		"start",
		fmt.Sprintf("--name=%s", p.Id),
		fmt.Sprintf("--addr=%s", net.JoinHostPort("", fmt.Sprint(p.Port))),
	}
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	attr.Env = os.Environ()
	bin, _ := exec.LookPath(os.Args[0])
	proc, err := os.StartProcess(bin, args, &attr)
	if err != nil {
		desc := fmt.Sprintf("os.StartProcess error:%s", err)
		return nil, status.NewStatusDesc(scode.PMCallGoLibError, desc)
	}
	p.Pid = int64(proc.Pid)
	m.DB.Put(BucketKeyPort, fmt.Sprint(p.Port), p.Id)
	m.DB.Put(BucketKeyPid, fmt.Sprint(p.Pid), p.Id)

	// 等待进程启动并同步信息
	time.Sleep(time.Millisecond)
	for i := 0; i < 10; i++ {
		err = m.updatePlayerInfo(p)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, status.UpdateStatus(err)
	}
	return p, nil
}

// 同步信息
func (m *PlayerManager) updatePlayerInfo(p *ipm.PlayerInfo) error {
	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(p.Port))
	c, err := iclient.NewPlayerClient(addr)
	if err != nil {
		return status.UpdateStatus(err)
	}
	defer c.Close()
	resp, err := c.SyncInfo(context.Background(), p)
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	p.Etag = resp.Etag
	return nil
}

// 退出进程
func (m *PlayerManager) SignOutPlayer(req *ipm.PlayerSignOutRequest) error {
	m.Lock()
	p, exist := m.players[req.Pid]
	m.Unlock()
	if !exist {
		desc := fmt.Sprintf("player %s is not exist", req.Pid)
		return status.NewStatusDesc(scode.PlayerNotExist, desc)
	}
	if p.Etag != req.Etag {
		desc := fmt.Sprintf("etag %s is wrong", req.Pid)
		return status.NewStatusDesc(scode.PlayerAuthWrong, desc)
	}
	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(p.Port))
	c, err := iclient.NewPlayerClient(addr)
	if err != nil {
		return status.UpdateStatus(err)
	}
	defer c.Close()
	resp, err := c.Stop(context.Background(), p)
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	if resp.Status.Code != 0 {
		return status.UpdateStatus(resp.Status)
	}
	return nil
}
