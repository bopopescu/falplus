package service

import (
	"api/ipm"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"os"
	"pm"
	"server"
	"status"
)

var (
	pmServer *PMService
	log      = logrus.WithFields(logrus.Fields{"pkg": "pm/service"})
)

type PMService struct {
	pm      *pm.PlayerManager
	service *server.Service
}

func NewPMServer(configFilePath, name, proto, addr string) *PMService {
	if pmServer == nil {
		pmServer = &PMService{}
	}
	service := server.NewService(configFilePath, "pm", name, proto, addr, pmServer)
	pmServer.service = service
	return pmServer
}

func (m *PMService) Init(c *goconfig.ConfigFile) {
	m.pm = pm.NewPlayerManager()
}

func (m *PMService) Signal(sig os.Signal) bool {
	log.Infof("PMService Signal Process... signal %s", sig.String())
	if sig.String() == "terminated" {

	}
	return true
}

func (m *PMService) Start() {
	ipm.RegisterPMServer(m.service.Server, m)
	reflection.Register(m.service.Server)
	m.service.Start()
}

func (m *PMService) PlayerCreate(ctx context.Context, req *ipm.PlayerCreateRequest) (*ipm.PlayerCreateResponse, error) {
	resp := &ipm.PlayerCreateResponse{}
	resp.Status = status.SuccessStatus
	p, err := m.pm.CreatePlayer(req)
	if err != nil {
		resp.Status = status.NewStatus(3000, err.Error())
		return resp, nil
	}
	resp.Player = p
	return resp, nil
}

func (m *PMService) PlayerDelete(ctx context.Context, req *ipm.PlayerDeleteRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	err := m.pm.DeletePlayer(req.Pid)
	if err != nil {
		resp.Status = status.NewStatus(3001, err.Error())
	}
	return resp, nil
}

func (m *PMService) PlayerList(ctx context.Context, req *ipm.PlayerListRequest) (*ipm.PlayerListResponse, error) {
	resp := &ipm.PlayerListResponse{}
	resp.Status = status.SuccessStatus
	for pid, info := range m.pm.GetAllPlayerInfo() {
		if pid == req.Pid || req.Pid == "" {
			resp.Players = append(resp.Players, info)
		}
	}
	return resp, nil
}

func (m *PMService) PlayerSignIn(ctx context.Context, req *ipm.PlayerSignInRequest) (*ipm.PlayerSignInResponse, error) {
	resp := &ipm.PlayerSignInResponse{}
	resp.Status = status.SuccessStatus
	p, err := m.pm.StartPlayer(req)
	if err != nil {
		resp.Status = status.NewStatus(3000, err.Error())
		return resp, nil
	}
	resp.Port = p.Port
	resp.Etag = p.Etag
	return resp, nil
}

func (m *PMService) PlayerSignOut(ctx context.Context, req *ipm.PlayerSignOutRequest) (*ipm.PMDefaultResponse, error) {
	resp := &ipm.PMDefaultResponse{}
	resp.Status = status.SuccessStatus
	return resp, nil
}
