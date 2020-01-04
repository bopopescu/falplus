package util

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/procfs"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
	"os"
	"strconv"
	"strings"
)

func MkdirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
			return err
		}
	}
	return nil
}

func ValidIP(ip string) bool {
	ret := net.ParseIP(ip)
	if ret != nil {
		return true
	}

	return false
}

func ValidPort(strPort string) bool {
	port, err := strconv.Atoi(strPort)
	if err != nil {
		return false
	}

	return port >= 0 && port <= 65535
}

func ValidIPAddr(addr string) bool {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	if host != "" && !ValidIP(host) {
		return false
	}

	if !ValidPort(port) {
		return false
	}

	return true
}

func PrintStructObject(data interface{}) {
	output, err := json.MarshalIndent(data, "", "\t")
	if err == nil {
		fmt.Println(string(output))
	} else {
		fmt.Println(err)
	}
}

func GetIPAddrFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if ok && p.Addr != nil {
		return p.Addr.String()
	}
	return ""
}

func GetIPv4Addr() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok {
				if p4 := ipnet.IP.To4(); len(p4) == net.IPv4len && !strings.HasPrefix(ipnet.IP.String(), "127"){
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

func InSliceString(val string, slice []string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func FindPids(match []string) ([]int, error) {
	procs, err := procfs.AllProcs()
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0)
	for _, proc := range procs {
		cmdLine, err := proc.CmdLine()
		if err != nil {
			continue
		}
		flag := true
		for _, m := range match {
			if !InSliceString(m, cmdLine) {
				flag = false
				break
			}
		}
		if flag {
			stat, err := proc.NewStat()
			if err != nil {
				return nil, err
			}
			pids = append(pids, stat.PID)
		}

	}
	return pids, nil
}