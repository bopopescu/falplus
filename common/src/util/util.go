package util

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
	"os"
	"strconv"
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
				if p4 := ipnet.IP.To4(); len(p4) == net.IPv4len {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}
