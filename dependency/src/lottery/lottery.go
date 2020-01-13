package main

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lotteryNum = 10
	prizeNum = []int{1,2,3}
	picType = "gif"
)

func GetIPv4Addr() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok {
				if p4 := ipnet.IP.To4(); len(p4) == net.IPv4len &&
					(strings.HasPrefix(ipnet.IP.String(), "192") || strings.HasPrefix(ipnet.IP.String(), "10")) {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

func main() {
	addr := net.JoinHostPort(GetIPv4Addr(), "12580")
	content := fmt.Sprintf("http://%s/lottery", addr)
	file := filepath.Join(filepath.Dir(os.Args[0]), "二维码.png")
	if err := qrcode.WriteFile(content, qrcode.Medium, 256, file); err != nil {
		fmt.Println(err)
		return
	}

	c := filepath.Join(filepath.Dir(os.Args[0]), "配置文件.txt")
	config(c)

	mux := http.NewServeMux()
	h := &myHandler{Nums: make([]int, lotteryNum+1)}
	conf := filepath.Join(filepath.Dir(os.Args[0]), "进度")
	data, err := ioutil.ReadFile(conf)
	if err == nil {
		json.Unmarshal(data, h)
	}
	mux.Handle("/lottery", h)

	log.Println("抽奖服务开启:")
	log.Fatal(http.ListenAndServe(addr, mux))
}

type myHandler struct{
	sync.Mutex
	Count int
	Nums  []int
}

func config(conf string) {
	cfg, err := goconfig.LoadConfigFile(conf)
	if err != nil {
		log.Fatalf("goconfig LoadConfigFile: %s, error: %s", conf, err.Error())
	}
	lotteryNum = cfg.MustInt("lottery", "lotteryNum", 1200)
	prizestr := cfg.MustValueArray("lottery", "prizeNum", ",")
	if len(prizestr) > 0 {
		prizeNum = make([]int, 0)
		for _, value := range prizestr {
			v, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("奖品等级数量错误")
				fmt.Println(err)
				os.Exit(-1)
			}
			prizeNum = append(prizeNum, v)
		}
	}
	picType = cfg.MustValue("lottery", "picType", "png")
}

func prize(i int) int {
	count := 0
	for index, value := range prizeNum {
		if i > count && i <= count+value {
			return index+1
		}
		count += value
	}
	return len(prizeNum)+1
}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	if m.Count >= lotteryNum {
		w.Write([]byte("抽奖已经结束"))
		return
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(lotteryNum-m.Count) + 1
	count := 0
	j := 1
	for ; j <= lotteryNum; j++ {
		if m.Nums[j] == 0 {
			count++
		}
		if count == index {
			m.Nums[j] = m.Count+1
			break
		}
	}
	m.Count++

	level := prize(j)
	if level == len(prizeNum)+1 {
		log.Printf("共计抽奖%d次，本次未中奖\n", m.Count)
	} else {
		log.Printf("共计抽奖%d次，本次抽中%d等奖\n", m.Count, level)
	}
	file := filepath.Join(filepath.Dir(os.Args[0]), fmt.Sprintf("%d.%s", level, picType))
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", picType))
	w.Write(data)

	conf := filepath.Join(filepath.Dir(os.Args[0]), "进度")
	data, err = json.Marshal(m)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	if err := ioutil.WriteFile(conf, data, 0666); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
