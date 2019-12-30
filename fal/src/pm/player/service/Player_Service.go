package service

import (
	"api/ipm"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"os"
	"pm/player"
	"scode"
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
	geted   bool
}

func NewPlayerServer(conf, name, proto, addr string) *PlayerService {
	srv := &PlayerService{
		player: player.NewPlayer(),
	}
	service := server.NewService(conf, "player", name, proto, addr, srv)
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

// 同步信息到PM
func (p *PlayerService) SyncInfo(ctx context.Context, req *ipm.PlayerInfo) (*ipm.PlayerInfo, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	p.player.UpdateInfo(req)
	return req, nil
}

// 退出进程
func (p *PlayerService) Stop(ctx context.Context, req *ipm.PlayerInfo) (*ipm.PMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.CloseConnGame(req.Etag)
	if err != nil {
		log.Error("CloseConnGame error:%s", err)
	}
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	return resp, nil
}

// 与game端建立数据连接
func (p *PlayerService) Attach(ctx context.Context, req *ipm.AttachRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.ConnGame(req.Etag, req.GamePort)
	if err != nil {
		log.Errorf("ConnGame error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 断开连接
func (p *PlayerService) Detach(ctx context.Context, req *ipm.DetachRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := p.player.CloseConnGame(req.Etag)
	if err != nil {
		log.Errorf("CloseConnGame error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 获取游戏消息（看牌）
func (p *PlayerService) GetMessage(ctx context.Context, req *ipm.GetMessageRequest) (*ipm.GetMessageResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.GetMessageResponse{}
	resp.Status = status.SuccessStatus
	msg, err := p.player.ShowGameMsg(req.Etag)
	if err != nil {
		log.Errorf("ShowGameMsg error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	resp.Gmsg = msg
	p.geted = true
	return resp, nil
}

// 出牌
func (p *PlayerService) PutMessage(ctx context.Context, req *ipm.PutMessageRequest) (*ipm.PMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if !p.geted {
		desc := fmt.Sprintf("must get before put")
		resp.Status = status.NewStatusDesc(scode.PutMessageBeforeGet, desc)
		return resp, nil
	}
	err := p.player.PutCards(req)
	if err != nil {
		log.Errorf("PutCards error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	p.geted = false
	return resp, nil
}
