package service

import (
	"api/igm"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"gm"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"os"
	"server"
	"status"
	"util"
)

var (
	gmServer *GMService
	log      = logrus.WithFields(logrus.Fields{"pkg": "gm/service"})
)

type GMService struct {
	gm      *gm.GameManager
	service *server.Service
}

func NewGMServer(conf, name, proto, addr string) *GMService {
	if gmServer == nil {
		gmServer = &GMService{}
	}
	service := server.NewService(conf, "gm", name, proto, addr, gmServer)
	gmServer.service = service
	return gmServer
}

func (m *GMService) Init(c *goconfig.ConfigFile) {
	m.gm = gm.NewGameManager()
	go func() {
		err := m.gm.InitUpdate()
		if err != nil {
			log.Error(status.UpdateStatus(err).Details())
		}
	}()
}

func (m *GMService) Signal(sig os.Signal) bool {
	log.Infof("GMService Signal Process... signal %s", sig.String())
	if sig.String() == "terminated" {
		m.gm.Stop()
	}
	return true
}

func (m *GMService) Start() {
	igm.RegisterGMServer(m.service.Server, m)
	reflection.Register(m.service.Server)
	m.service.Start()
}

func (m *GMService) GameCreate(ctx context.Context, req *igm.GameCreateRequest) (*igm.GameCreateResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GameCreateResponse{}
	resp.Status = status.SuccessStatus
	g, err := m.gm.CreateGame(req)
	if err != nil {
		resp.Status = status.NewStatus(2000, err.Error())
		return resp, nil
	}
	resp.Gid = g.Gid
	resp.Port = g.Port
	resp.GameType = g.GameType
	return resp, nil
}

func (m *GMService) GameDelete(ctx context.Context, req *igm.GameDeleteRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := m.gm.DeleteGame(req.Gid)
	if err != nil {
		resp.Status = status.NewStatus(2001, err.Error())
	}
	return resp, nil
}

func (m *GMService) GameList(ctx context.Context, req *igm.GameListRequest) (*igm.GameListResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GameListResponse{}
	resp.Status = status.SuccessStatus
	for gid, info := range m.gm.GetAllGameInfo() {
		if gid == req.Gid || req.Gid == "" {
			resp.Games = append(resp.Games, info)
		}
	}
	return resp, nil
}

func (m *GMService) GameStart(ctx context.Context, req *igm.GameStartRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}

func (m *GMService) GameStop(ctx context.Context, req *igm.GameStopRequest) (*igm.GMDefaultResponse, error) {
	log.Infof("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}
