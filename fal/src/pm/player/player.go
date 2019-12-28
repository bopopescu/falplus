package player

import (
	"api/igm"
	"api/ipm"
	"card"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"iclient"
	"sync"
)

type Player struct {
	Id      string
	User    string
	Pwd     string
	Score   int
	Cards   []int
	gClient *iclient.GameClient
	stream  igm.Game_PlayerConnClient
	close   chan struct{}
	msgChan chan *message
	resp    chan *igm.PlayerMessage
	req     *igm.GameMessage
	Etag    string
	port    int64
}

var (
	log    = logrus.WithFields(logrus.Fields{"pkg": "pm/player/service"})
	player *Player
)

func NewPlayer() *Player {
	if player == nil {
		player = &Player{
			close:   make(chan struct{}),
			msgChan: make(chan *message, 1),
			resp:    make(chan *igm.PlayerMessage),
		}
	}
	return player
}

func (p *Player) Close() {
	if p.gClient != nil {
		p.gClient.Close()
	}
}

func (p *Player) UpdateInfo(info *ipm.PlayerInfo) {
	info.Etag = uuid.New().String()
	p.Etag = info.Etag
	p.User = info.Name
	p.Pwd = info.Password
	p.port = info.Port
	p.Id = info.Id
}

func (p *Player) ShowGameMsg(etag string) (*igm.GameMessage, error) {
	if p.Etag != etag {
		return nil, fmt.Errorf("auth error")
	}
	if len(p.msgChan) == 0 {
		return nil, fmt.Errorf("last req is not finish")
	}
	msg := <-p.msgChan
	go func() {
		msg.pMsg = <-p.resp
		msg.Done()
	}()
	p.req = msg.gMsg
	return msg.gMsg, msg.err
}

func (p *Player) PutCards(req *ipm.PutMessageRequest) error {
	if p.Etag != req.Etag {
		return fmt.Errorf("auth error")
	}
	if !p.checkRespValid(req.Pmsg) {
		return fmt.Errorf("the resp is invalid")
	}
	p.resp <- req.Pmsg
	return nil
}

func (p *Player) checkRespValid(resp *igm.PlayerMessage) bool {
	if p.req.RoundOwner == p.Id {
		// 争夺地主中只能选择抢或者过
		// 出牌判断是否可出
		switch p.req.MsgType {
		case igm.Get:
			if resp.MsgType != igm.Get && resp.MsgType != igm.Pass {
				return false
			}
			resp.PutCards = nil
		case igm.Put:
			lastRepeated, lastValue, lastLength := card.GetRepeatNumAndValue(p.req.LastCards)
			lastType := card.GetCardsType(lastRepeated, p.req.LastCards)

			repeated, value, length := card.GetRepeatNumAndValue(resp.PutCards)
			curType := card.GetCardsType(repeated, resp.PutCards)

			// 牌类型正确且上次牌无人要
			if curType != card.ValueTypeUnknown && p.req.LastId == p.Id {
				return true
			}

			// 牌类型未知 或者非炸弹的情况下类型不等 或者类型相等长度不等 或者类型相等值较小
			if curType == card.ValueTypeUnknown ||
				curType != lastType && curType != card.ValueTypeBomb ||
				curType == lastType && length != lastLength ||
				curType == lastType && value < lastValue {
				return false
			}
		}
	} else {
		// 无牌权无需关心
		resp.MsgType = igm.Msg
		resp.PutCards = nil
	}

	return true
}

func (p *Player) ConnGame(etag, addr string) error {
	if p.Etag != etag {
		return fmt.Errorf("auth error")
	}
	c, err := iclient.NewGameClient(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	stream, err := c.PlayerConn(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}
	pMsg := &igm.PlayerMessage{
		MsgType:  igm.Msg,
		PlayerId: p.Id,
	}
	err = stream.Send(pMsg)
	if err != nil {
		log.Error(err)
		return err
	}
	p.gClient = c
	p.stream = stream
	go p.ListenGame()
	return nil
}

func (p *Player) CloseConnGame(etag string) error {
	if p.Etag != etag {
		return fmt.Errorf("auth error")
	}
	if p.gClient == nil {
		return fmt.Errorf("conn is nil")
	}
	return p.gClient.Close()
}

type message struct {
	gMsg *igm.GameMessage
	pMsg *igm.PlayerMessage
	err  error
	sync.WaitGroup
}

func (p *Player) ListenGame() {
	for {
		msg := &message{}
		req, err := p.stream.Recv()
		if err != nil {
			msg.err = err
		}
		msg.gMsg = req
		msg.Add(1)
		select {
		case <-p.close:
			return
		case p.msgChan <- msg:
		}
		msg.Wait()
		err = p.stream.Send(msg.pMsg)
		if err != nil {
			panic(err)
		}
	}
}
