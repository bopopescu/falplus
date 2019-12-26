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
)

var log = logrus.WithFields(logrus.Fields{"pkg": "gm/game/service"})

type GameService struct {
	game *game.Game
	service *server.Service
	gameOver chan struct{}
}

func NewGameServer(configFilePath, name, proto, addr string) *GameService {
	srv := &GameService{gameOver:make(chan struct{})}
	service := server.NewService(configFilePath, "game", name, proto, addr, srv)
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

func (m *GameService) Start() {
	igm.RegisterGameServer(m.service.Server, m)
	reflection.Register(m.service.Server)
	m.service.Start()
}

func (m *GameService) PlayerConn(stream igm.Game_PlayerConnServer) error {
	if len(m.game.Players) < 3 {
		m.game.Players = append(m.game.Players, stream)
	} else {
		return fmt.Errorf("the room if full, please select anothor")
	}
	select {
	case <-m.gameOver:
		return nil
	}
}

func (m *GameService) GameStart (ctx context.Context, req *igm.GameStartRequest) (*igm.GMDefaultResponse, error) {
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}

func (m *GameService) GameStop (ctx context.Context, req *igm.GameStopRequest) (*igm.GMDefaultResponse, error) {
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}
