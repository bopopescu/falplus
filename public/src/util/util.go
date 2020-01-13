package util

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"io/ioutil"
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

func AllPids() ([]int, error) {
	var p []int
	d, err := os.Open("/proc")
	if err != nil {
		return p, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return p, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	for _, n := range names {
		pid, err := strconv.Atoi(n)
		if err != nil {
			continue
		}
		p = append(p, pid)
	}

	return p, nil
}

func FindPids(match []string) ([]int, error) {
	allPids, err := AllPids()
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, pid := range allPids {
		if matchPid(pid, match) {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

func matchPid(pid int, match []string) bool {
	data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false
	}
	//[:len(data)-1]
	cmdLine := strings.Split(string(data), string(byte(0)))
	cmdLineMap := make(map[string]struct{})
	for _, m := range cmdLine {
		cmdLineMap[m] = struct{}{}
	}
	for _, m := range match {
		if _, ok := cmdLineMap[m]; !ok {
			return false
		}
	}
	return true
}
