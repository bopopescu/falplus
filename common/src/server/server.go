package server

import (
	"api/idebug"
	"debugservice"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/goconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"logger"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"util"
)

const (
	LOGROTATE_TEMPLATE = `{
        size 500M
        rotate 10
        missingok
        compress
        delaycompress
        copytruncate
        notifempty
        create 0644 root root
        postrotate
                invoke-rc.d rsyslog rotate > /dev/null
        endscript
}`

	LOGROTATE_PATH = "/etc/logrotate.d/"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "service"})
)

type Handler interface {
	Init(c *goconfig.ConfigFile)
	Signal(sig os.Signal) bool //return true to stop daemon
}

type Service struct {
	*ServiceConfig
	handler      Handler
	tcpListener  net.Listener
	unixListener net.Listener
	Server       *grpc.Server
}

type ServiceConfig struct {
	Root  string
	Name  string
	Proto string
	Addr  string
	Sock  string
}

func NewService(conf, module, name, proto, addr string, handler Handler) *Service {
	//cfg, err := goconfig.LoadConfigFile(conf)
	//if err != nil {
	//	log.Fatalf("goconfig LoadConfigFile: %s, error: %s", conf, err.Error())
	//}

	rootDir := "/var/log/fal"
	if err := util.MkdirIfNotExists(rootDir); err != nil {
		log.Fatalf("MkdirIfNotExists %s error: %s", rootDir, err.Error())
	}

	sconfig := &ServiceConfig{}
	sconfig.Root = filepath.Join(rootDir, module)
	if err := util.MkdirIfNotExists(sconfig.Root); err != nil {
		log.Fatalf("MkdirIfNotExists %s error: %s", sconfig.Root, err.Error())
	}

	if len(name) != 0 {
		sconfig.Name = name
	} else {
		sconfig.Name = module
	}

	if len(proto) == 0 {
		sconfig.Proto = "tcp"
	} else {
		sconfig.Proto = proto
	}

	if sconfig.Proto == "tcp" {
		if util.ValidIPAddr(addr) {
			sconfig.Addr = addr
		} else {
			sconfig.Addr = ":9999"
		}
	}

	//初始化日志
	logLevel := "debug"
	logrus.SetLevel(transferLevel(logLevel))
	logdir := sconfig.Root

	logrotatefile := filepath.Join(LOGROTATE_PATH, sconfig.Name)
	file, err := os.OpenFile(logrotatefile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("OpenFile %s error: %s", logrotatefile, err.Error())
	}

	logName := fmt.Sprintf("%s.log", sconfig.Name)
	logPath := filepath.Join(logdir, logName)
	if module == "game" {
		logPath = filepath.Join(logdir, "game", logName)
	}
	if module == "player" {
		logPath = filepath.Join(logdir, "player", logName)
	}
	if err := util.MkdirIfNotExists(filepath.Dir(logPath)); err != nil {
		log.Fatalf("MkdirIfNotExists %s error: %s", sconfig.Root, err.Error())
	}
	file.WriteString(logPath)
	file.WriteString("\n")
	file.WriteString(LOGROTATE_TEMPLATE)
	file.WriteString("\n")
	file.Close()

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("OpenFile %s error: %s", logPath, err.Error())
	}
	syscall.Dup3(int(logFile.Fd()), 1, 0)
	syscall.Dup3(int(logFile.Fd()), 2, 0)
	logrus.SetFormatter(&logger.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000000"})
	logrus.SetOutput(logFile)
	log.Infof("module %s logLevel %s ", module, logLevel)

	srv := &Service{
		handler:       handler,
		ServiceConfig: sconfig,
	}

	//设置信号处理
	sigs := make(chan os.Signal, 5)

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 102400)
			length := runtime.Stack(buf, true)
			log.Errorf("Service main goroutine panic, error: %s \n stack:\n%s\n", err, string(buf[:length]))
			os.Exit(-1)
		}
	}()

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, os.Kill, syscall.SIGTERM)
	go func() {
		for {
			sig := <-sigs
			log.Infof("Caught signal %s, trying to shutdown service", sig)
			if srv.handler.Signal(sig) {
				srv.Server.Stop()
				//buf := make([]byte, 102400)
				//length := runtime.Stack(buf, true)
				//log.Infof("Service main goroutine panic, error: %s \n stack:\n%s\n", err, string(buf[:length]))
				break
			}
		}
	}()

	if srv.Proto == "unix" {
		if err := util.MkdirIfNotExists(filepath.Dir(srv.Sock)); err != nil {
			log.Fatalf("MkdirIfNotExists %s error: %s", srv.Sock, err.Error())
		}

		if _, err := os.Stat(srv.Sock); err == nil {
			log.Warnf("Remove previous sockfile at %s", srv.Sock)
			if err = os.Remove(srv.Sock); err != nil {
				log.Fatalf("Remove previous sock file %s error:%s", srv.Sock, err.Error())
			}
		}

		srv.unixListener, err = net.Listen(srv.Proto, srv.Sock)
		if err != nil {
			log.Fatalf("Unix listen error: %s", err.Error())
		}
	}

	if srv.Proto == "tcp" {
		srv.tcpListener, err = net.Listen(srv.Proto, srv.Addr)
		if err != nil {
			log.Fatalf("TCP listen addr %s error: %s", srv.Addr, err.Error())
		}
	}

	var kaepOpt = grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	})
	var kaspOpt = grpc.KeepaliveParams(keepalive.ServerParameters{
		//MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		//MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		//MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  60 * time.Second,  // Ping the client if it is idle for 10 seconds to ensure the connection is still active
		Timeout:               30 * time.Second,  // Wait 5 second for the ping ack before assuming the connection is dead
	})
	srv.Server = grpc.NewServer(kaepOpt, kaspOpt)

	srv.handler.Init(&goconfig.ConfigFile{})

	idebug.RegisterDebugServer(srv.Server, &debugservice.DebugService{})

	return srv
}

func (srv *Service) Start() {
	if srv.unixListener != nil {
		srv.Server.Serve(srv.unixListener)
	}

	if srv.tcpListener != nil {
		srv.Server.Serve(srv.tcpListener)
	}
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
