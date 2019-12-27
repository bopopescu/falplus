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
	Players []igm.Game_PlayerConnServer
	Cards   [][]int64
	LastId  string
	msgChan chan *message
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
		case msg := <-g.msgChan:
			g.sendMessage(msg)
			msg.Done()
		}
	}
}

// 向所有玩家发送信息，获取牌权持有者的响应
func (g *Game) sendMessage(msg *message) {
	for index, stream := range g.Players {
		msg.gMsg.LastCards = g.Cards[3]     // 场上牌
		msg.gMsg.YourCards = g.Cards[index] // 玩家手里牌
		msg.gMsg.LastId = g.LastId
		if msg.owner == index {
			msg.gMsg.RoundOwner = true // 牌权
		}
		err := stream.Send(msg.gMsg)
		if err != nil {
			msg.err = err
			return
		}
		r, err := stream.Recv()
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
func (g *Game) qiangdizhu() (int, error) {
	for index := 0; index < 3; index++ {
		resp, err := g.round(igm.QiangDiZhu, index, nil, "")
		if err != nil {
			return 0, err
		}
		if resp.MsgType == igm.QiangDiZhu {
			g.Cards[index] = append(g.Cards[index], g.Cards[3]...)
			g.Cards[3] = g.Cards[3][:0]
			return index, nil
		}
	}
	return 0, fmt.Errorf("nobody get landlord")
}

// 回合指定回合类型，牌权，当前场上牌，场上牌所有者
func (g *Game) round(rtype int64, owner int, pcard []int64, powner string) (*igm.PlayerMessage, error) {
	msg := &message{}
	gMsg := &igm.GameMessage{
		MsgType:   rtype,
		LastId:    powner,
		LastCards: pcard,
	}
	msg.gMsg = gMsg
	msg.owner = owner
	msg.Add(1)
	g.msgChan <- msg
	msg.Wait()
	return msg.pMsg, msg.err
}

// 指定牌权，当前回合类型，场上牌，场上牌所属。
func (g *Game) gameLogical() {
	lIndex, err := g.qiangdizhu()
	if err != nil {
		panic(err)
	}
	msg := &message{owner: lIndex}
	gMsg := &igm.GameMessage{MsgType: igm.ChuPai}
	msg.gMsg = gMsg
	for {
		msg.Add(1)
		g.msgChan <- msg
		msg.Wait()
		if msg.err != nil {
			panic(err)
		}
		resp := msg.pMsg
		if resp.MsgType == igm.Pass {
			continue
		}
		// 更新并判断
		if g.updateCards(msg.owner, resp.PutCards) {

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
