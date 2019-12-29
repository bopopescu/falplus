package debugservice

import (
	"api/idebug"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"scode"
	"status"
	"time"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "debugservice"})
	lis net.Listener
)

type DebugService struct {
}

func (ds *DebugService) Prof(ctx context.Context, req *idebug.ProfRequest) (*idebug.ProfResponse, error) {
	resp := &idebug.ProfResponse{}
	resp.Status = status.SuccessStatus

	f, err := os.Create(req.Path)
	if err != nil {
		log.Errorf("Create path %s error %s", req.Path, err.Error())
		desc := fmt.Sprintf("Create path %s error %s", req.Path, err)
		return nil, status.NewStatusDesc(scode.DebugCallGoLibError, desc)
	}
	if req.Name == "cpu" {
		pprof.StartCPUProfile(f)
		time.Sleep(time.Second * time.Duration(req.Time))
		pprof.StopCPUProfile()
	} else {
		pprof.Lookup(req.Name).WriteTo(f, 0)
	}
	f.Close()
	return resp, nil
}

func (ds *DebugService) Stats(ctx context.Context, req *idebug.StatsRequest) (*idebug.StatsResponse, error) {
	resp := &idebug.StatsResponse{}
	resp.Status = status.SuccessStatus

	if req.Name == "routine" {
		resp.Data = []byte(fmt.Sprintf("routine %d\n", runtime.NumGoroutine()))
		return resp, nil
	}

	if req.Name == "stack" {
		buf := make([]byte, 409600)
		length := runtime.Stack(buf, true)
		resp.Data = buf[:length]
		return resp, nil
	}

	if req.Name == "memstats" {
		memstats := &runtime.MemStats{}
		runtime.ReadMemStats(memstats)
		mem, err := json.MarshalIndent(memstats, "", "\t")
		if err != nil {
			mem = []byte(err.Error())
		}
		resp.Data = mem
		return resp, nil
	}
	return resp, nil
}

func (ds *DebugService) GetLogLevel(ctx context.Context, req *google_protobuf.Empty) (*idebug.LogResponse, error) {
	resp := &idebug.LogResponse{}
	resp.Status = status.SuccessStatus
	resp.Level = logrus.GetLevel().String()
	return resp, nil
}

func (ds *DebugService) SetLogLevel(ctx context.Context, req *idebug.LogRequest) (*idebug.LogResponse, error) {
	resp := &idebug.LogResponse{}
	resp.Status = status.SuccessStatus
	logrus.SetLevel(transferLevel(req.Level))
	resp.Level = logrus.GetLevel().String()
	return resp, nil
}

func transferLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

func (ds *DebugService) StartPprof(ctx context.Context, req *idebug.StartPprofRequest) (*idebug.StartPprofResponse, error) {
	resp := &idebug.StartPprofResponse{}
	resp.Status = status.SuccessStatus
	listen, err := net.Listen("tcp", req.Addr)
	if err != nil {
		log.Errorf("start debug service with addr:%s, error:%s", req.Addr, err)
		desc := fmt.Sprintf("start debug service with addr:%s, error:%s", req.Addr, err)
		resp.Status = status.NewStatusDesc(scode.DebugCallGoLibError, desc)
		return resp, nil
	}
	lis = listen
	log.Infof("start debug pprof with addr:%s", lis.Addr())
	resp.Listenaddr = lis.Addr().String()
	go http.Serve(lis, nil)
	return resp, nil
}

func (ds *DebugService) StopPprof(ctx context.Context, req *google_protobuf.Empty) (*idebug.StopPprofResponse, error) {
	resp := &idebug.StopPprofResponse{}
	resp.Status = status.SuccessStatus
	if lis == nil {
		log.Errorf("close debug service error:pprof server is not exist")
		desc := fmt.Sprintf("close debug service error:pprof server is not exist")
		resp.Status = status.NewStatusDesc(scode.DebugAlreadyStop, desc)
		return resp, nil
	}
	if err := lis.Close(); err != nil {
		log.Errorf("close debug service error:%s", err)
		desc := fmt.Sprintf("close debug service error:%s", err)
		resp.Status = status.NewStatus(scode.DebugCallGoLibError, desc)
		return resp, nil
	}
	lis = nil
	return resp, nil
}
