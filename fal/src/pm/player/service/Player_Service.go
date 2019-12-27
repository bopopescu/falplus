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

func (p *Player_Service) SyncInfo(ctx context.Context, req *ipm.PlayerInfo) (*ipm.PlayerInfo, error) {
	p.player.UpdateInfo(req)
	return req, nil
}

func (p *Player_Service) Attach(ctx context.Context, req *ipm.AttachRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.ConnGame(req.Etag, req.GamePort)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
	}
	return resp, nil
}

func (p *Player_Service) Detach(ctx context.Context, req *ipm.DetachRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.CloseConnGame(req.Etag)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
	}
	return resp, nil
}

func (p *Player_Service) GetMessage(ctx context.Context, req *ipm.GetMessageRequest) (*ipm.GetMessageResponse, error) {
	resp := &ipm.GetMessageResponse{}
	resp.Status = status.SuccessStatus
	msg, err := p.player.ShowGameMsg(req.Etag)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
		return resp, nil
	}
	resp.Gmsg = msg
	return resp, nil
}

func (p *Player_Service) PutMessage(ctx context.Context, req *ipm.PutMessageRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.PutCards(req)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
	}
	return resp, nil
}
