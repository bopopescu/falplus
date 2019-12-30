package service

import (
	"api/igm"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"gm"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"iclient"
	"net"
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
		if err := m.gm.InitUpdate(); err != nil {
			log.Errorf("InitUpdate error:%s", err)
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

// 创建并启动游戏
func (m *GMService) GameCreate(ctx context.Context, req *igm.GameCreateRequest) (*igm.GameCreateResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GameCreateResponse{}
	resp.Status = status.SuccessStatus
	g, err := m.gm.CreateGame(req)
	if err != nil {
		log.Error(err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	resp.Gid = g.Gid
	resp.Port = g.Port
	resp.GameType = g.GameType
	return resp, nil
}

// 停止并删除游戏
func (m *GMService) GameDelete(ctx context.Context, req *igm.GameDeleteRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.Exit(ctx, &igm.GameExitRequest{})
	}); err != nil {
		log.Errorf("Exit error:%s", err)
	}
	if err := m.gm.DeleteGame(req.Gid); err != nil {
		log.Error(err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 游戏列表
func (m *GMService) GameList(ctx context.Context, req *igm.GameListRequest) (*igm.GameListResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GameListResponse{}
	resp.Status = status.SuccessStatus
	for gid, info := range m.gm.GetAllGameInfo() {
		if gid == req.Gid || req.Gid == "" {
			resp.Games = append(resp.Games, info)
		}
	}
	return resp, nil
}

// 玩家加入房间
func (m *GMService) GameAddPlayer(ctx context.Context, req *igm.AddPlayerRequest) (*igm.AddPlayerResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.AddPlayerResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.AddPlayer(ctx, req)
	}); err != nil {
		log.Errorf("AddPlayer error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	port := m.gm.GetGameInfo(req.Gid).Port
	resp.GameAddr = net.JoinHostPort(util.GetIPv4Addr(), fmt.Sprint(port))
	return resp, nil
}

// 玩家离开房间
func (m *GMService) GameDelPlayer(ctx context.Context, req *igm.DelPlayerRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.DelPlayer(ctx, req)
	}); err != nil {
		log.Errorf("DelPlayer error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 开始游戏
func (m *GMService) GameStart(ctx context.Context, req *igm.GameStartRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.Start(ctx, req)
	}); err != nil {
		log.Error("start error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 停止游戏
func (m *GMService) GameStop(ctx context.Context, req *igm.GameStopRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.Stop(ctx, req)
	}); err != nil {
		log.Errorf("Stop error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}

// 退出游戏进程
func (m *GMService) GameExit(ctx context.Context, req *igm.GameExitRequest) (*igm.GMDefaultResponse, error) {
	log.Debugf("get client addr %s request:%v", util.GetIPAddrFromCtx(ctx), req)
	resp := &igm.GMDefaultResponse{}
	resp.Status = status.SuccessStatus
	if err := m.gm.DefaultGameResponse(req.Gid, func(c *iclient.GameClient) (*igm.GMDefaultResponse, error) {
		return c.Exit(ctx, req)
	}); err != nil {
		log.Errorf("Exit error:%s", err)
		resp.Status = status.UpdateStatus(err)
		return resp, nil
	}
	return resp, nil
}
