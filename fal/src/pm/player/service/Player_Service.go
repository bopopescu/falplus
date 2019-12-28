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
	"time"
	"util"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "pm/player/service"})
)

type PlayerService struct {
	player  *player.Player
	service *server.Service
}

func NewPlayerServer(conf, name, proto, addr string) *PlayerService {
	srv := &PlayerService{
		player: player.NewPlayer(),
	}
	service := server.NewService(conf, "", name, proto, addr, srv)
	srv.service = service
	return srv
}

func (p *PlayerService) Init(c *goconfig.ConfigFile) {

}

func (p *PlayerService) Signal(sig os.Signal) bool {
	log.Infof("PlayerService Signal Process... signal %s", sig.String())
	if sig.String() == "terminated" {

	}
	return true
}

func (p *PlayerService) Start() {
	ipm.RegisterPlayerServer(p.service.Server, p)
	reflection.Register(p.service.Server)
	p.service.Start()
}

func (p *PlayerService) SyncInfo(ctx context.Context, req *ipm.PlayerInfo) (*ipm.PlayerInfo, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	p.player.UpdateInfo(req)
	return req, nil
}

func (p *PlayerService) Stop(ctx context.Context, req *ipm.PlayerInfo) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if req.Etag != p.player.Etag {
		resp.Status = status.NewStatus(3001, "etag error")
		log.Errorf("etag error")
		return resp, nil
	}
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	return resp, nil
}

func (p *PlayerService) Attach(ctx context.Context, req *ipm.AttachRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.ConnGame(req.Etag, req.GamePort)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
		log.Error(err)
	}
	return resp, nil
}

func (p *PlayerService) Detach(ctx context.Context, req *ipm.DetachRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.CloseConnGame(req.Etag)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
		log.Error(err)
	}
	return resp, nil
}

func (p *PlayerService) GetMessage(ctx context.Context, req *ipm.GetMessageRequest) (*ipm.GetMessageResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.GetMessageResponse{}
	resp.Status = status.SuccessStatus
	msg, err := p.player.ShowGameMsg(req.Etag)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
		log.Error(err)
		return resp, nil
	}
	resp.Gmsg = msg
	return resp, nil
}

func (p *PlayerService) PutMessage(ctx context.Context, req *ipm.PutMessageRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.PutCards(req)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
		log.Error(err)
	}
	return resp, nil
}
