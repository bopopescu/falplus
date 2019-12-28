package game

import (
	"api/igm"
	"card"
	"fmt"
	"sync"
)

type Game struct {
	sync.RWMutex
	Status  int64
	Id      string
	Port    int64
	Pid     int
	Players []*PInfo
	Cards   [][]int64
	LastId  string
	LastWin int
	msgChan chan *message
	Stop    chan struct{}
}

type PInfo struct {
	Stream igm.Game_PlayerConnServer
	Id     string
	Addr   string
}

type message struct {
	owner int
	gMsg  *igm.GameMessage
	pMsg  *igm.PlayerMessage
	err   error
	sync.WaitGroup
}

var game *Game

const (
	Init  = 0
	Start = 1
)

func NewGame() *Game {
	game = &Game{
		msgChan: make(chan *message),
		Stop:    make(chan struct{}),
	}
	return game
}

func (g *Game) Start() {
	g.Lock()
	g.Status = Start
	cards := card.DistributeCards(54)
	g.Cards = append(g.Cards, cards[:13], cards[13:26], cards[26:39], cards[52:])
	g.Unlock()
	for {
		select {
		case <-g.Stop:
			return
		case msg := <-g.msgChan:
			g.sendMessage(msg)
			msg.Done()
		}
	}
}

// 向所有玩家发送信息，获取牌权持有者的响应
func (g *Game) sendMessage(msg *message) {
	for index, p := range g.Players {
		msg.gMsg.LastCards = g.Cards[3]     // 场上牌
		msg.gMsg.YourCards = g.Cards[index] // 玩家手里牌
		msg.gMsg.LastId = g.LastId
		if msg.owner == index {
			msg.gMsg.RoundOwner = g.Players[index].Id // 牌权
		}
		err := p.Stream.Send(msg.gMsg)
		if err != nil {
			msg.err = err
			return
		}
		r, err := p.Stream.Recv()
		if err != nil {
			msg.err = err
			return
		}
		if msg.owner == index {
			msg.pMsg = r
		}
	}
}

// 争夺地主
func (g *Game) fightForLandlord(cur int) (int, error) {
	for i := cur; i < cur+3; i++ {
		index := i % 3
		resp, err := g.round(igm.Get, index)
		if err != nil {
			return 0, err
		}
		if resp.MsgType == igm.Get {
			g.Cards[index] = append(g.Cards[index], g.Cards[3]...)
			g.Cards[3] = g.Cards[3][:0]
			return index, nil
		}
	}
	return 0, fmt.Errorf("nobody get landlord")
}

// 回合指定回合类型，牌权
func (g *Game) round(rtype int64, owner int) (*igm.PlayerMessage, error) {
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
func (g *Game) GameLogical() {
	lIndex, err := g.fightForLandlord(g.LastWin)
	if err != nil {
		panic(err)
	}
	for {
		cur := lIndex
		resp, err := g.round(igm.Put, cur)
		if err != nil {
			panic(err)
		}
		lIndex++
		lIndex %= 3
		if resp.MsgType == igm.Pass {
			continue
		}
		g.Cards[3] = resp.PutCards
		g.LastId = g.Players[cur].Id
		// 更新并判断
		if g.updateCards(cur, resp.PutCards) {
			g.LastWin = cur
			close(g.Stop)
			return
		}
	}
}

// 更新玩家手牌并判断是否结束游戏
func (g *Game) updateCards(index int, cards []int64) bool {
	tmp := make(map[int64]struct{})
	for _, seq := range cards {
		tmp[seq] = struct{}{}
	}
	var c []int64
	for _, seq := range g.Cards[index] {
		if _, ok := tmp[seq]; !ok {
			c = append(c, seq)
		}
	}
	if len(c) == 0 {
		return true
	}
	g.Cards[index] = c
	return false
}
