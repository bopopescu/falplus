package pm

import (
	"api/ipm"
	"fdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"iclient"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
	"util"
)

const (
	PMDB             = "/var/lib/fal/pm.db"
	BucketPlayerList = "players"
	PlayerPrefix     = "player-"
	BucketKeyPid     = "pid"
	BucketKeyName    = "name"
	BucketKeyPwd     = "password"
)

var (
	pm  *PlayerManager
	log = logrus.WithFields(logrus.Fields{"pkg": "pm"})
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
		Id:       id,
		Name:     req.Name,
		Password: req.Password,
	}
	kvs := make(map[string]string)
	kvs[BucketKeyPid] = id
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

func (m *PlayerManager) StartPlayer(req *ipm.PlayerSignInRequest) (*ipm.PlayerInfo, error) {
	m.Lock()
	defer m.Unlock()
	p, exist := m.players[req.Pid]
	if !exist {
		log.Errorf("player %s is not create", req.Pid)
		return p, fmt.Errorf("")
	}
	if p.Name != req.Name || p.Password != req.Password {
		log.Errorf("player %s is name %s or password %s error", req.Pid, req.Name, req.Password)
		return p, fmt.Errorf("")
	}
	p.Port = m.assignPort(p.Port)

	args := []string{
		filepath.Base(os.Args[0]),
		"player",
		"start",
		"--name", p.Id,
		"--addr", net.JoinHostPort("", fmt.Sprint(p.Port)),
	}
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	attr.Env = os.Environ()
	bin, _ := exec.LookPath(os.Args[0])
	proc, err := os.StartProcess(bin, args, &attr)
	if err != nil {
		fmt.Println(err)
		return p, err
	}
	p.Pid = int64(proc.Pid)
	time.Sleep(time.Millisecond)
	for i := 0; i < 10; i++ {
		err = m.updatePlayerInfo(p)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return p, err
	}
	return p, nil
}

func (m *PlayerManager) updatePlayerInfo(p *ipm.PlayerInfo) error {
	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(p.Port))
	c, err := iclient.NewPlayerClient(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer c.Close()
	resp, err := c.SyncInfo(context.Background(), p)
	if err != nil {
		log.Error(err)
		return err
	}
	p.Etag = resp.Etag
	return nil
}

func (m *PlayerManager) SignOutPlayer(req *ipm.PlayerSignOutRequest) error {
	p, exist := m.players[req.Pid]
	if !exist {
		return fmt.Errorf("pid %s is not exist", req.Pid)
	}
	if p.Etag != req.Etag {
		return fmt.Errorf("etag %s is wrong", req.Pid)
	}
	addr := net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(p.Port))
	c, err := iclient.NewPlayerClient(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer c.Close()
	resp, err := c.Stop(context.Background(), p)
	if err != nil {
		log.Error(err)
		return err
	}
	if resp.Status.Code != 0 {
		log.Error(resp.Status)
		return resp.Status
	}
	return nil
}
