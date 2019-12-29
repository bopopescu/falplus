package service

import (
	"api/igm"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"gm/game/fal"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"os"
	"server"
	"status"
	"sync"
	"time"
	"util"
)

var log = logrus.WithFields(logrus.Fields{"pkg": "gm/game/service"})

type GameService struct {
	game     igm.Game
	service  *server.Service
	operate  *sync.Map
	gameOver chan struct{}
}

func NewGameServer(configFilePath, name, proto, addr string) *GameService {
	srv := &GameService{
		operate:  &sync.Map{},
		gameOver: make(chan struct{}),
	}
	service := server.NewService(configFilePath, "", name, proto, addr, srv)
	srv.service = service
	return srv
}

func (m *GameService) Init(c *goconfig.ConfigFile) {
	m.game = fal.NewGame()
}

func (m *GameService) Signal(sig os.Signal) bool {
	log.Debugf("GameService Signal Process... signal %s", sig.String())
	if sig.String() == "terminated" {

	}
	return true
}

func (m *GameService) Run() {
	igm.RegisterGameServer(m.service.Server, m)
	reflection.Register(m.service.Server)
	m.service.Start()
}

// 玩家加入房间
func (m *GameService) AddPlayer(ctx context.Context, req *igm.AddPlayerRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.game.AddPlayer(req.Pid); err != nil {
		log.Errorf("AddPlayer error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 玩家离开房间
func (m *GameService) DelPlayer(ctx context.Context, req *igm.DelPlayerRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.game.DelPlayer(req.Pid); err != nil {
		log.Errorf("DelPlayer error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 玩家建立数据连接
func (m *GameService) PlayerConn(stream igm.Game_PlayerConnServer) error {
	log.Debugf("get client addr %s", util.GetIPAddrFromCtx(stream.Context()))
	leave, err := m.game.PlayerConn(stream)
	if err != nil {
		log.Errorf("PlayerConn error:%s", err)
		return status.UpdateStatus(err)
	}
	select {
	case <-leave:
	case <-m.gameOver:
	}
	return nil
}

// 房主开始游戏
func (m *GameService) Start(ctx context.Context, req *igm.GameStartRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.game.Start(req.Pid); err != nil {
		log.Errorf("Start error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 房主结束游戏
func (m *GameService) Stop(ctx context.Context, req *igm.GameStopRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.game.Stop(req.Pid); err != nil {
		log.Errorf("Stop error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 关闭游戏进程
func (m *GameService) Exit(ctx context.Context, req *igm.GameExitRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	// todo 退出权限控制

	if err := m.game.Stop(req.Pid); err != nil {
		log.Errorf("Stop error:%s", err)
	}
	close(m.gameOver)
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	return resp, nil
}
