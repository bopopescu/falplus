package gm

import (
	"api/igm"
	"encoding/json"
	"fdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

func (m *GameManager) InitUpdate() error {
	m.Lock()
	defer m.Unlock()

	if m.isInit {
		panic("InitUpdate USM twice")
	}
	m.isInit = true

	return m.gameInit()
}

func (m *GameManager) gameInit() error {
	buckets, err := m.DB.GetAllBucket()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(buckets)
	for _, b := range buckets {
		kvs, err := m.DB.GetAllKV(b)
		if err != nil {
			log.Error(err)
			return err
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

func (m *GameManager) CreateGame(req *igm.GameCreateRequest) (*igm.GameInfo, error) {
	m.Lock()
	defer m.Unlock()
	id := req.Gid
	if req.Gid != "" {
		_, exist := m.games[req.Gid]
		if exist {
			return nil, fmt.Errorf("game already exist")
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
		log.Error(err)
		return nil, err
	}

	kvs := make(map[string]string)
	kvs[BucketKeyGid] = id
	kvs[BucketKeyType] = fmt.Sprint(g.GameType)
	kvs[BucketKeyState] = fmt.Sprint(g.State)
	kvs[BucketKeyPort] = fmt.Sprint(g.Port)
	kvs[BucketKeyPid] = fmt.Sprint(g.Pid)
	data := map[string]map[string]string{id: kvs}
	m.games[id] = g
	err = m.DB.Put(id, fmt.Sprint(req.GameType), BucketGameList)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return g, m.DB.PutBatch(data)
}

func (m *GameManager) assignPort(port int64) int64 {
	if port != 0 {
		return port
	}
	lis, err := net.Listen("tcp", ":")
	if err != nil {

	}
	lis.Close()
	return int64(lis.Addr().(*net.TCPAddr).Port)
}

func (m *GameManager) DeleteGame(id string) error {
	m.Lock()
	defer m.Unlock()
	_, exist := m.games[id]
	if !exist {
		return fmt.Errorf("game not exist")
	}
	delete(m.games, id)

	if err := m.DB.Delete(id, BucketGameList); err != nil {
		log.Error(err)
	}
	if err := m.DB.DeleteBucket(id); err != nil {
		log.Error(err)
	}
	return nil
}

func (m *GameManager) GetAllGameInfo() map[string]*igm.GameInfo {
	games := make(map[string]*igm.GameInfo)
	for k, g := range m.games {
		games[k] = g
	}
	return games
}

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
	bin, _ := exec.LookPath(os.Args[0])
	proc, err := os.StartProcess(bin, args, &attr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	g.Pid = int64(proc.Pid)
	return nil
}

func (m *GameManager) stopGame(id string) error {
	return nil
}
