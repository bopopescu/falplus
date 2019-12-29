package fal

import (
	"api/igm"
	"card"
	"fmt"
	"github.com/Sirupsen/logrus"
	"scode"
	"status"
	"sync"
)

var log = logrus.WithFields(logrus.Fields{"pkg": "gm/game/fal"})

type GameFal struct {
	sync.RWMutex
	players    []*PInfo
	cards      [][]int64
	lastId     string
	lastWin    int
	msgChan    chan *message
	stopNormal chan struct{}
	stopForce  chan struct{}
	state      int
}

type PInfo struct {
	Stream igm.Game_PlayerConnServer
	Id     string
	Leave  chan struct{}
}

type message struct {
	owner int
	gMsg  *igm.GameMessage
	pMsg  *igm.PlayerMessage
	err   error
	sync.WaitGroup
}

var game *GameFal

func NewGame() *GameFal {
	game = &GameFal{
		msgChan:    make(chan *message),
		stopNormal: make(chan struct{}),
		stopForce:  make(chan struct{}),
	}
	return game
}

func (g *GameFal) AddPlayer(pid string) error {
	g.Lock()
	defer g.Unlock()
	if len(g.players) >= 3 {
		desc := fmt.Sprintf("game room is full")
		return status.NewStatusDesc(scode.GamePlayerIsFull, desc)
	}
	p := &PInfo{
		Id:    pid,
		Leave: make(chan struct{}),
	}
	g.players = append(g.players, p)
	return nil
}

func (g *GameFal) DelPlayer(pid string) error {
	g.Lock()
	defer g.Unlock()
	var players []*PInfo
	for _, p := range g.players {
		if p.Id != pid {
			players = append(players, p)
		} else {
			close(p.Leave)
		}
	}
	// 长度相等说明pid不存在
	if len(players) == len(g.players) {
		desc := fmt.Sprintf("player %s is not in this game", pid)
		return status.NewStatusDesc(scode.GamePlayerIsFull, desc)
	}
	g.players = players
	return nil
}

func (g *GameFal) PlayerConn(stream igm.Game_PlayerConnServer) (<-chan struct{}, error) {
	g.Lock()
	defer g.Unlock()
	pMsg, err := stream.Recv()
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return nil, status.NewStatusDesc(scode.GRPCError, desc)
	}
	var p *PInfo
	exist := false
	for _, p = range g.players {
		if p.Id == pMsg.Pid {
			exist = true
			p.Stream = stream
		}
	}
	if !exist {
		desc := fmt.Sprintf("player %s is not belong this game", pMsg.Pid)
		return nil, status.NewStatusDesc(scode.PlayerNotInTheGame, desc)
	}
	return p.Leave, nil
}

// 判断发起游戏开始玩家是否为房主,人数是否够三人（房主为第一进房间的人或者上次的赢家）
func (g *GameFal) Start(pid string) error {
	g.Lock()
	defer g.Unlock()
	if g.state == igm.Start {
		return status.NewStatus(scode.GameAlreadyStart)
	}
	g.state = igm.Start
	if len(g.players) < 3 {
		return status.NewStatus(scode.GamePlayerNotEnough)
	}
	if g.players[g.lastWin].Id != pid {
		desc := fmt.Sprintf("player %s is not host in this game", pid)
		return status.NewStatusDesc(scode.GamePlayerIsNotHost, desc)
	}
	go g.goStart()
	go g.goGameLogical()

	return nil
}

func (g *GameFal) Stop(pid string) error {
	g.Lock()
	defer g.Unlock()
	if g.state == igm.Stop {
		return status.NewStatus(scode.GameAlreadyStop)
	}
	g.state = igm.Stop
	if g.players[g.lastWin].Id != pid {
		desc := fmt.Sprintf("player %s is not host in this game", pid)
		return status.NewStatusDesc(scode.GamePlayerIsNotHost, desc)
	}
	close(g.stopForce)
	return nil
}

func (g *GameFal) State() int {
	return g.state
}

func (g *GameFal) goStart() {
	cards := card.DistributeCards(54)
	// 此方法赋值会复用cards的内存，造成意想不到的bug
	// g.Cards = append(g.Cards, cards[:17], cards[17:34], cards[34:51], cards[51:])
	g.cards = make([][]int64, 3)
	for index := range g.cards {
		g.cards[index] = make([]int64, 17)
		copy(g.cards[index], cards[index*17:(index+1)*17])
	}
	g.cards = append(g.cards, cards[51:])
	for {
		select {
		case <-g.stopNormal:
			return
		case msg := <-g.msgChan:
			g.sendMessage(msg)
			msg.Done()
		}
	}
}

// 向所有玩家发送信息，获取牌权持有者的响应
func (g *GameFal) sendMessage(msg *message) {
	wg := &sync.WaitGroup{}
	for index, p := range g.players {
		wg.Add(1)
		go func(index int, p *PInfo) {
			defer wg.Done()
			gMsg := *msg.gMsg
			gMsg.LastCards = g.cards[3]     // 场上牌
			gMsg.YourCards = g.cards[index] // 玩家手里牌
			gMsg.LastId = g.lastId
			gMsg.RoundOwner = g.players[msg.owner].Id // 牌权
			err := p.Stream.Send(&gMsg)
			if err != nil {
				desc := fmt.Sprintf("player %s GRPC Send error:%s", p.Id, err)
				msg.err = status.NewStatusDesc(scode.GRPCError, desc)
				return
			}
			r, err := p.Stream.Recv()
			if err != nil {
				desc := fmt.Sprintf("player %s GRPC Recv error:%s", p.Id, err)
				msg.err = status.NewStatusDesc(scode.GRPCError, desc)
				return
			}
			if msg.owner == index {
				msg.pMsg = r
			}
		}(index, p)
	}
	wg.Wait()
}

// 争夺地主
func (g *GameFal) fightForLandlord(cur int) (int, error) {
	for i := cur; i < cur+3; i++ {
		index := i % 3
		resp, err := g.round(igm.Get, index)
		if err != nil {
			return 0, status.UpdateStatus(err)
		}
		if resp.MsgType == igm.Get {
			g.cards[index] = append(g.cards[index], g.cards[3]...)
			g.cards[3] = g.cards[3][:0]
			return index, nil
		}
	}
	desc := fmt.Sprintf("nobody want get landlord")
	return 0, status.NewStatusDesc(scode.NobodyWantGetLord, desc)
}

// 回合指定回合类型，牌权
func (g *GameFal) round(rtype int64, owner int) (*igm.PlayerMessage, error) {
	msg := &message{}
	gMsg := &igm.GameMessage{MsgType: rtype}
	msg.gMsg = gMsg
	msg.owner = owner
	msg.Add(1)
	g.msgChan <- msg
	msg.Wait()
	return msg.pMsg, msg.err
}

// 指定牌权，当前回合类型，场上牌，场上牌所属。
func (g *GameFal) goGameLogical() {
	defer func() {
		close(g.stopNormal)
		g.state = igm.Stop
	}()
	lIndex, err := g.fightForLandlord(g.lastWin)
	if err != nil {
		if r := status.FromError(err); r.Code == scode.NobodyWantGetLord {
			g.round(igm.Gpass, g.lastWin)
		}
		log.Errorf("fightForLandlord error:%s", err)
		return
	}
	for {
		select {
		case <-g.stopForce:
			return
		default:
			cur := lIndex
			resp, err := g.round(igm.Put, cur)
			if err != nil {
				log.Errorf("round error:%s", err)
				return
			}
			// 牌权轮换，回合制
			lIndex++
			lIndex %= 3
			if resp.MsgType == igm.Pass {
				continue
			}
			g.cards[3] = resp.PutCards
			g.lastId = g.players[cur].Id
			// 更新并判断
			if g.updateCards(cur, resp.PutCards) {
				g.round(igm.Over, cur)
				g.lastWin = cur
				return
			}
		}
	}
}

// 更新玩家手牌并判断是否结束游戏
func (g *GameFal) updateCards(index int, cards []int64) bool {
	tmp := make(map[int64]struct{})
	for _, seq := range cards {
		tmp[seq] = struct{}{}
	}
	var c []int64
	for _, seq := range g.cards[index] {
		if _, ok := tmp[seq]; !ok {
			c = append(c, seq)
		}
	}
	if len(c) == 0 {
		return true
	}
	g.cards[index] = c
	return false
}
