package pm

import (
	"api/ipm"
	"fdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"net"
	"sync"
)

const (
	PMDB             = "/var/lib/fal/pm.db"
	BucketPlayerList = "players"
	PlayerPrefix     = "player-"
	BucketKeyPid     = "pid"
	BucketKeyName    = "name"
	BucketKeyPwd     = "password"
	BucketKeyPort    = "port"
)

var (
	pm  *PlayerManager
	log = logrus.WithFields(logrus.Fields{"pkg": "gm"})
)

type PlayerManager struct {
	sync.RWMutex
	players map[string]*ipm.PlayerInfo
	DB      *fdb.FalDB
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

func (m *PlayerManager) AddPlayer(user, pwd string) error {
	pass, err := m.DB.Get(user, BucketPlayerList)
	if err != nil || pass != pwd {
		return fmt.Errorf("account or sercret error")
	}
	err = m.DB.CreateBucket(user)
	if err != nil {
		panic(err)
	}
	m.Lock()
	defer m.Unlock()
	_, exist := m.players[user]
	if exist {
		return fmt.Errorf("player already exist")
	}
	m.players[user] = &ipm.PlayerInfo{}
	return nil
}

func (m *PlayerManager) CreatePlayer(req *ipm.PlayerCreateRequest) (*ipm.PlayerInfo, error) {
	m.Lock()
	defer m.Unlock()
	id := req.Pid
	if req.Pid != "" {
		_, exist := m.players[req.Pid]
		if exist {
			return nil, fmt.Errorf("game already exist")
		}
	} else {
		id = uuid.New().String()
	}

	id = PlayerPrefix + id

	p := &ipm.PlayerInfo{
		Pid:      id,
		Name:     req.Name,
		Password: req.Password,
		Port:     m.assignPort(req.Port),
	}
	kvs := make(map[string]string)
	kvs[BucketKeyPid] = id
	kvs[BucketKeyPort] = fmt.Sprint(p.Port)
	kvs[BucketKeyName] = p.Name
	kvs[BucketKeyPwd] = p.Password
	data := map[string]map[string]string{id: kvs}
	m.players[id] = p

	if err := m.DB.Put(id, fmt.Sprint(req.Name), BucketPlayerList); err != nil {
		log.Error(err)
		return nil, err
	}
	return p, m.DB.PutBatch(data)
}

func (m *PlayerManager) assignPort(port int64) int64 {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		lis, _ = net.Listen("tcp", ":")
	}
	lis.Close()
	return int64(lis.Addr().(*net.TCPAddr).Port)
}

func (m *PlayerManager) DeletePlayer(id string) error {
	m.Lock()
	defer m.Unlock()
	_, exist := m.players[id]
	if !exist {
		return fmt.Errorf("game not exist")
	}
	delete(m.players, id)

	if err := m.DB.Delete(id, BucketPlayerList); err != nil {
		log.Error(err)
	}
	if err := m.DB.DeleteBucket(id); err != nil {
		log.Error(err)
	}
	return nil
}

func (m *PlayerManager) GetAllPlayerInfo() map[string]*ipm.PlayerInfo {
	players := make(map[string]*ipm.PlayerInfo)
	for k, g := range m.players {
		players[k] = g
	}
	return players
}
