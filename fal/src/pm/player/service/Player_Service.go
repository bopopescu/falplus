package service

import (
	"api/ipm"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"os"
	"pm/player"
	"server"
	"status"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "gm/service"})
)

type Player_Service struct {
	player  *player.Player
	service *server.Service
}

func NewPlayerServer(conf, name, proto, addr string) *Player_Service {
	srv := &Player_Service{
		player: player.NewPlayer(),
	}
	service := server.NewService(conf, "player", name, proto, addr, srv)
	srv.service = service
	return srv
}

func (m *Player_Service) Init(c *goconfig.ConfigFile) {

}

func (m *Player_Service) Signal(sig os.Signal) bool {
	log.Infof("Player_Service Signal Process... signal %s", sig.String())
	if sig.String() == "terminated" {

	}
	return true
}

func (p *Player_Service) Start() {
	ipm.RegisterPlayerServer(p.service.Server, p)
	reflection.Register(p.service.Server)
	p.service.Start()
}

func (p *Player_Service) SignIn(ctx context.Context, req *ipm.SignInRequest) (*ipm.SignInResponse, error) {
	resp := &ipm.SignInResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}

func (p *Player_Service) SignOut(ctx context.Context, req *ipm.SignOutRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}

func (p *Player_Service) GetMessage(ctx context.Context, req *ipm.GetMessageRequest) (*ipm.GetMessageResponse, error) {
	resp := &ipm.GetMessageResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}

func (p *Player_Service) PutMessage(ctx context.Context, req *ipm.PutMessageRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}
