package service

import (
	"api/igm"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"gm/game"
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
	game     *game.Game
	service  *server.Service
	gameOver chan struct{}
	operate  *sync.Map
}

func NewGameServer(configFilePath, name, proto, addr string) *GameService {
	srv := &GameService{
		gameOver: make(chan struct{}),
		operate:  &sync.Map{},
	}
	service := server.NewService(configFilePath, "", name, proto, addr, srv)
	srv.service = service
	return srv
}

func (m *GameService) Init(c *goconfig.ConfigFile) {
	m.game = game.NewGame()
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
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if len(m.game.Players) >= 3 {
		resp.Status = status.NewStatus(2001, "game room is full")
		log.Errorf("game room is full")
		return resp, nil
	}
	p := &game.PInfo{
		Id:   req.PlayerId,
		Addr: req.PlayerAddr,
	}
	m.game.Players = append(m.game.Players, p)
	return resp, nil
}

// 玩家建立数据连接
func (m *GameService) PlayerConn(stream igm.Game_PlayerConnServer) error {
	log.Infof("get client addr %s", util.GetIPAddrFromCtx(stream.Context()))
	pMsg, err := stream.Recv()
	if err != nil {
		log.Errorf("recv pmsg error:%s", err.Error())
		return fmt.Errorf("recv pmsg error:%s", err.Error())
	}
	exist := false
	for _, p := range m.game.Players {
		if p.Id == pMsg.PlayerId {
			exist = true
			p.Stream = stream
		}
	}
	if !exist {
		log.Errorf("player is not belong this game")
		return fmt.Errorf("player is not belong this game")
	}
	select {
	case <-m.gameOver:
		return nil
	}
}

// 房主开始游戏
func (m *GameService) Start(ctx context.Context, req *igm.GameStartRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	// 判断发起游戏开始玩家是否为房主,人数是否够三人（房主为第一进房间的人或者上次的赢家）
	exist := false
	count := 0
	for index, p := range m.game.Players {
		if p.Id == req.Pid && m.game.LastWin == index {
			exist = true
		}
		if p.Stream != nil {
			count++
		}
	}
	if !exist || count != 3 {
		resp.Status = status.NewStatus(2001, "you can not start game")
		log.Errorf("you can not start game")
		return resp, nil
	}

	// 判断游戏是否已经开始
	lock := &sync.Mutex{}
	_, loaded := m.operate.LoadOrStore("GameStart", lock)
	if loaded {
		resp.Status = status.NewStatus(2001, "game already start")
		log.Errorf("game already start")
		return resp, nil
	}
	go m.game.Start()
	go m.game.GameLogical()
	return resp, nil
}

// 房主结束游戏
func (m *GameService) Stop(ctx context.Context, req *igm.GameStopRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	exist := false
	for index, p := range m.game.Players {
		if p.Id == req.Pid && m.game.LastWin == index {
			exist = true
		}
	}
	if !exist {
		resp.Status = status.NewStatus(2001, "you can not stop game")
		log.Errorf("you can not stop game")
		return resp, nil
	}
	close(m.game.Stop)
	m.operate.Delete("GameStart")
	return resp, nil
}

// 关闭游戏进程
func (m *GameService) Exit(ctx context.Context, req *igm.GameExitRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	//exist := false
	//for index, p := range m.game.Players {
	//	if p.Id == req.PlayerId && m.game.LastWin == index {
	//		exist = true
	//	}
	//}
	//if !exist {
	//	resp.Status = status.NewStatus(2001, "you can not exit game")
	//	return resp, nil
	//}
	close(m.game.Stop)
	m.operate.Delete("GameStart")
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	return resp, nil
}
