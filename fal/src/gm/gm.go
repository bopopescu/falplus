package gm

import (
	"api/igm"
	"encoding/json"
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

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "gm"})
)

const (
	GMDB             = "/var/lib/fal/gm.db"
	BucketGameList   = "games"
	GamePrefix       = "game-"
	BucketKeyPid     = "pid"
	BucketKeyType    = "type"
	BucketKeyGid     = "gid"
	BucketKeyState   = "state"
	BucketKeyPlayers = "players"
	BucketKeyPort    = "port"
)

var gm *GameManager
var once = new(sync.Once)

type GameManager struct {
	sync.RWMutex
	games  map[string]*igm.GameInfo
	DB     *fdb.FalDB
	isInit bool
}

func NewGameManager() *GameManager {
	once.Do(func() {
		// 切记放在init，否则将会导致一直卡在opendb
		db := fdb.NewDB(GMDB)
		err := db.CreateBucket(BucketGameList)
		if err != nil {
			panic(err)
		}
		gm = &GameManager{
			games: make(map[string]*igm.GameInfo),
			DB:    db,
		}
	})
	return gm
}

func (m *GameManager) Stop() {
	if m.DB != nil {
		m.DB.Close()
	}
}

// 初始化
func (m *GameManager) InitUpdate() error {
	m.Lock()
	defer m.Unlock()

	if m.isInit {
		panic("InitUpdate GM twice")
	}
	m.isInit = true

	if err := m.gameInit(); err != nil {
		log.Error(err)
	}

	return m.scanGameInfo()
}

func (m *GameManager) scanGameInfo() error {
	match := []string{"fal", "game", "start"}
	pids, err := util.FindPids(match)
	if err != nil {
		return err
	}

	for _, pid := range pids {
		file := fmt.Sprintf("/proc/%d/cmdline", pid)
		output, err := ioutil.ReadFile(file)
		if err != nil {
			log.Errorf("read game info error:%s", err.Error())
			continue
		}
		cmdLine := strings.Split(string(output[:len(output)-1]), string(byte(0)))
		var gid string
		var addr string
		for _, str := range cmdLine {
			if strings.HasPrefix(str, "--name=") {
				gid = strings.TrimPrefix(str, "--name=")
			}
			if strings.HasPrefix(str, "--addr=") {
				addr = strings.TrimPrefix(str, "--addr=")
			}
		}

		if _, exist := m.games[gid]; exist {
			_, portstr, _ := net.SplitHostPort(addr)
			m.games[gid].Pid = int64(pid)
			m.games[gid].Port, _ = strconv.ParseInt(portstr, 10, 64)
			if err := m.updateGameInfo(m.games[gid]); err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

// 从数据库读入游戏数据
func (m *GameManager) gameInit() error {
	buckets, err := m.DB.GetAllBucket()
	if err != nil {
		desc := fmt.Sprintf("GetAllBucket error:%s", err)
		return status.NewStatusDesc(scode.GMDBOperateError, desc)
	}
	log.Debug(buckets)
	for _, b := range buckets {
		if strings.HasPrefix(b, GamePrefix) {

			kvs, err := m.DB.GetAllKV(b)
			if err != nil {
				desc := fmt.Sprintf("GetAllKV error:%s", err)
				return status.NewStatusDesc(scode.GMDBOperateError, desc)
			}

			port, _ := strconv.ParseInt(kvs[BucketKeyPort], 10, 64)
			state, _ := strconv.ParseInt(kvs[BucketKeyState], 10, 64)
			pid, _ := strconv.ParseInt(kvs[BucketKeyPid], 10, 64)
			gType, _ := strconv.ParseInt(kvs[BucketKeyType], 10, 64)
			players := make(map[string]string)
			json.Unmarshal([]byte(kvs[BucketKeyPlayers]), &players)
			info := &igm.GameInfo{
				Gid:      kvs[BucketKeyGid],
				Port:     port,
				State:    state,
				Pid:      pid,
				GameType: gType,
				Players:  players,
			}
			m.games[info.Gid] = info
		}
	}
	return nil
}

// 创建游戏并保存在数据库
func (m *GameManager) CreateGame(req *igm.GameCreateRequest) (*igm.GameInfo, error) {
	m.Lock()
	defer m.Unlock()
	id := req.Gid
	if req.Gid != "" {
		_, exist := m.games[req.Gid]
		if exist {
			desc := fmt.Sprintf("game %s already exist", req.Gid)
			return nil, status.NewStatusDesc(scode.GameAlreadyExist, desc)
		}
	} else {
		id = uuid.New().String()
	}
	id = GamePrefix + id

	g := &igm.GameInfo{
		Gid:      id,
		GameType: req.GameType,
		Port:     m.assignPort(req.Port),
		State:    igm.Exit,
	}

	m.games[id] = g
	var err error
	defer func() {
		if err != nil {
			delete(m.games, id)
		}
	}()
	err = m.startGame(g)
	if err != nil {
		return nil, status.UpdateStatus(err)
	}

	kvs := make(map[string]string)
	kvs[BucketKeyGid] = id
	kvs[BucketKeyType] = fmt.Sprint(g.GameType)
	kvs[BucketKeyState] = fmt.Sprint(g.State)
	kvs[BucketKeyPort] = fmt.Sprint(g.Port)
	kvs[BucketKeyPid] = fmt.Sprint(g.Pid)
	players, _ := json.Marshal(g.Players)
	kvs[BucketKeyPlayers] = string(players)
	data := map[string]map[string]string{id: kvs}

	// 读写文件放开锁保证性能
	m.Unlock()
	defer m.Lock()
	if err = m.DB.Put(id, fmt.Sprint(req.GameType), BucketGameList); err != nil {
		desc := fmt.Sprintf("put key:%s value:%s bucket:%s error:%s", id, fmt.Sprint(req.GameType), BucketGameList, err)
		return nil, status.NewStatusDesc(scode.GMDBOperateError, desc)
	}
	if err = m.DB.PutBatch(data); err != nil {
		desc := fmt.Sprintf("PutBatch data:%v error:%s", data, err)
		return nil, status.NewStatusDesc(scode.GMDBOperateError, desc)
	}
	return g, nil
}

func (m *GameManager) UpdateGameInfo(g *igm.GameInfo) error {
	m.Lock()
	defer m.Unlock()
	return m.updateGameInfo(g)
}

// 同步信息
func (m *GameManager) updateGameInfo(g *igm.GameInfo) error {
	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(g.Port))
	c, err := iclient.NewGameClient(addr)
	if err != nil {
		return status.UpdateStatus(err)
	}
	defer c.Close()
	resp, err := c.SyncInfo(context.Background(), &igm.GameInfo{})
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	g.State = resp.State
	g.Players = resp.Players
	return nil
}

// 分配端口
func (m *GameManager) assignPort(port int64) int64 {
	lis, err := net.Listen("tcp", net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(port)))
	if err != nil {
		log.Info(err)
		lis, err = net.Listen("tcp", ":")
	}
	if err != nil {
		panic("assignPort err")
	}
	lis.Close()
	return int64(lis.Addr().(*net.TCPAddr).Port)
}

// 删除游戏
func (m *GameManager) DeleteGame(id string) error {
	m.Lock()
	defer m.Unlock()
	_, exist := m.games[id]
	if !exist {
		desc := fmt.Sprintf("game %s not exist", id)
		return status.NewStatusDesc(scode.GameNotExist, desc)
	}
	delete(m.games, id)

	m.Unlock()
	defer m.Lock()
	if err := m.DB.Delete(id, BucketGameList); err != nil {
		log.Error(err)
	}
	if err := m.DB.DeleteBucket(id); err != nil {
		log.Error(err)
	}
	return nil
}

// 获取所有游戏信息
func (m *GameManager) GetAllGameInfo() map[string]*igm.GameInfo {
	m.RLock()
	defer m.RUnlock()
	games := make(map[string]*igm.GameInfo)
	for k, g := range m.games {
		games[k] = g
	}
	return games
}

// 获取指定游戏信息
func (m *GameManager) GetGameInfo(gid string) *igm.GameInfo {
	m.RLock()
	defer m.RUnlock()
	g, _ := m.games[gid]
	return g
}

// 设置游戏状态
func (m *GameManager) SetState(gid string, state int64) {
	m.Lock()
	defer m.Unlock()
	m.games[gid].State = state
}

func (m *GameManager) StartGame(g *igm.GameInfo) error {
	m.Lock()
	defer m.Unlock()
	return m.startGame(g)
}

// 启动游戏进程
func (m *GameManager) startGame(g *igm.GameInfo) error {
	port := m.assignPort(g.Port)
	args := []string{
		filepath.Base(os.Args[0]),
		"game",
		"start",
		fmt.Sprintf("--name=%s", g.Gid),
		fmt.Sprintf("--addr=%s", net.JoinHostPort("", fmt.Sprint(port))),
	}

	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	attr.Env = os.Environ()
	bin, err := exec.LookPath(os.Args[0])
	if err != nil {
		desc := fmt.Sprintf("exec.LookPath %s error:%s", os.Args[0], err)
		return status.NewStatusDesc(scode.GMCallGoLibError, desc)
	}

	proc, err := os.StartProcess(bin, args, &attr)
	if err != nil {
		desc := fmt.Sprintf("os.StartProcess %s %v %v error:%s", bin, args, attr, err)
		return status.NewStatusDesc(scode.GMCallGoLibError, desc)
	}

	m.games[g.Gid].Pid = int64(proc.Pid)
	m.games[g.Gid].State = igm.Running
	if port != g.Port {
		m.games[g.Gid].Port = port
		m.DB.Put(BucketKeyPort, fmt.Sprint(port), g.Gid)
	}

	// 等待进程启动并同步信息
	time.Sleep(time.Millisecond)
	for i := 0; i < 10; i++ {
		err = m.updateGameInfo(m.games[g.Gid])
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return status.UpdateStatus(err)
	}

	return nil
}

// 默认响应格式的服务请求
func (m *GameManager) DefaultGameResponse(gid string, f func(c *iclient.GameClient) (*igm.GMDefaultResponse, error)) error {
	m.RLock()
	g, exist := m.games[gid]
	m.RUnlock()
	if !exist {
		desc := fmt.Sprintf("game %s not exist", gid)
		return status.NewStatusDesc(scode.GameNotExist, desc)
	}

	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(g.Port))
	c, err := iclient.NewGameClient(addr)
	if err != nil {
		return status.UpdateStatus(err)
	}
	defer c.Close()
	resp, err := f(c)
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	if resp.Status.Code != 0 {
		return status.UpdateStatus(resp.Status)
	}

	return nil
}

// 发送game请求 reqOpr发送请求 respOpr处理响应
func (m *GameManager) ConcurrenceGameOperate(gid string, reqOpr func(*igm.GameInfo) interface{}, respOpr func(interface{})) error {
	allGame := m.GetAllGameInfo()
	count := 0
	mc := make(chan interface{}, len(allGame))

	for _, game := range allGame {
		if gid == "all" || game.Gid == gid {
			count++
			go func(g *igm.GameInfo) {
				mc <- reqOpr(g)
			}(game)
		}
	}

	if count == 0 {
		desc := fmt.Sprintf("game %s not exist", gid)
		return status.NewStatusDesc(scode.GameNotExist, desc)
	}

	for count > 0 {
		r := <-mc
		respOpr(r)
		count--
	}

	return nil
}
