package gm

import (
	"api/igm"
	"encoding/json"
	"fdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"iclient"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"scode"
	"status"
	"strconv"
	"strings"
	"sync"
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

type GameManager struct {
	sync.RWMutex
	games  map[string]*igm.GameInfo
	DB     *fdb.FalDB
	isInit bool
}

func NewGameManager() *GameManager {
	db := fdb.NewDB(GMDB)
	err := db.CreateBucket(BucketGameList)
	if err != nil {
		panic(err)
	}
	gm = &GameManager{
		games: make(map[string]*igm.GameInfo),
		DB:    db,
	}
	return gm
}

func (m *GameManager) Stop() {

}

// 初始化
func (m *GameManager) InitUpdate() error {
	m.Lock()
	defer m.Unlock()

	if m.isInit {
		panic("InitUpdate USM twice")
	}
	m.isInit = true

	return m.gameInit()
}

// 从数据库读入游戏数据
func (m *GameManager) gameInit() error {
	buckets, err := m.DB.GetAllBucket()
	if err != nil {
		desc := fmt.Sprintf("GetAllBucket error:%s", err)
		return status.NewStatusDesc(scode.GMDBOperateError, desc)
	}
	log.Info(buckets)
	for _, b := range buckets {
		kvs, err := m.DB.GetAllKV(b)
		if err != nil {
			desc := fmt.Sprintf("GetAllKV error:%s", err)
			return status.NewStatusDesc(scode.GMDBOperateError, desc)
		}
		if strings.HasPrefix(b, GamePrefix) {
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
	}

	err := m.startGame(g)
	if err != nil {
		return nil, status.UpdateStatus(err)
	}

	kvs := make(map[string]string)
	kvs[BucketKeyGid] = id
	kvs[BucketKeyType] = fmt.Sprint(g.GameType)
	kvs[BucketKeyState] = fmt.Sprint(g.State)
	kvs[BucketKeyPort] = fmt.Sprint(g.Port)
	kvs[BucketKeyPid] = fmt.Sprint(g.Pid)
	data := map[string]map[string]string{id: kvs}
	m.games[id] = g

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

// 分配端口
func (m *GameManager) assignPort(port int64) int64 {
	lis, err := net.Listen("tcp", net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(port)))
	if err != nil {
		log.Error(err)
		lis, _ = net.Listen("tcp", ":")
	}
	lis.Close()
	return int64(lis.Addr().(*net.TCPAddr).Port)
}

// 删除游戏
func (m *GameManager) DeleteGame(id string) error {
	m.Lock()
	_, exist := m.games[id]
	if !exist {
		desc := fmt.Sprintf("game %s not exist", id)
		return status.NewStatusDesc(scode.GameNotExist, desc)
	}
	delete(m.games, id)
	m.Unlock()

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

// 启动游戏进程
func (m *GameManager) startGame(g *igm.GameInfo) error {
	args := []string{
		filepath.Base(os.Args[0]),
		"game",
		"start",
		"--name", g.Gid,
		"--addr", net.JoinHostPort("", fmt.Sprint(g.Port)),
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
		desc := fmt.Sprintf("os.StartProcess %s %v %v error:%s", bin, args, attr)
		return status.NewStatusDesc(scode.GMCallGoLibError, desc)
	}
	g.Pid = int64(proc.Pid)
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
