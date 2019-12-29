package iclient

import (
	"api/igm"
	"api/ipm"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"scode"
	"status"
	"time"
)

// 加入超时机制解决由于网络断开导致客户端一直处于连接状态不关闭
var kpOpt = grpc.WithKeepaliveParams(keepalive.ClientParameters{
	Time:                60 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             30 * time.Second, // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
})

type GMClient struct {
	igm.GMClient
	conn *grpc.ClientConn
}

type GameClient struct {
	igm.GameClient
	conn *grpc.ClientConn
}

type PMClient struct {
	ipm.PMClient
	conn *grpc.ClientConn
}

type PlayerClient struct {
	ipm.PlayerClient
	conn *grpc.ClientConn
}

func NewGMClient(addr string) (*GMClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), kpOpt)
	if err != nil {
		desc := fmt.Sprintf("Dial addr %s error %s", addr, err)
		return nil, status.NewStatusDesc(scode.GRPCError, desc)
	}
	return &GMClient{
		GMClient: igm.NewGMClient(conn),
		conn:     conn,
	}, nil
}

func NewGameClient(addr string) (*GameClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), kpOpt)
	if err != nil {
		desc := fmt.Sprintf("Dial addr %s error %s", addr, err)
		return nil, status.NewStatusDesc(scode.GRPCError, desc)
	}
	return &GameClient{
		GameClient: igm.NewGameClient(conn),
		conn:       conn,
	}, nil
}

func NewPMClient(addr string) (*PMClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), kpOpt)
	if err != nil {
		desc := fmt.Sprintf("Dial addr %s error %s", addr, err)
		return nil, status.NewStatusDesc(scode.GRPCError, desc)
	}
	return &PMClient{
		PMClient: ipm.NewPMClient(conn),
		conn:     conn,
	}, nil
}

func NewPlayerClient(addr string) (*PlayerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), kpOpt)
	if err != nil {
		desc := fmt.Sprintf("Dial addr %s error %s", addr, err)
		return nil, status.NewStatusDesc(scode.GRPCError, desc)
	}
	return &PlayerClient{
		PlayerClient: ipm.NewPlayerClient(conn),
		conn:         conn,
	}, nil
}

func (c *GMClient) Close()error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *GameClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *PMClient) Close()error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *PlayerClient) Close()error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
