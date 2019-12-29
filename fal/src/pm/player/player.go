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
	"scode"
	"status"
	"sync"
)

var (
	log    = logrus.WithFields(logrus.Fields{"pkg": "pm/player/service"})
	player *Player
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
	connErr error
}

type message struct {
	gMsg *igm.GameMessage
	pMsg *igm.PlayerMessage
	err  error
	sync.WaitGroup
}

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
		close(p.close)
	}
}

// 同步信息到GM
func (p *Player) UpdateInfo(info *ipm.PlayerInfo) {
	info.Etag = uuid.New().String()
	p.Etag = info.Etag
	p.User = info.Name
	p.Pwd = info.Password
	p.port = info.Port
	p.Id = info.Id
}

// 获取game端推送的消息
func (p *Player) ShowGameMsg(etag string) (*igm.GameMessage, error) {
	// 收发信息出错
	if p.connErr != nil {
		return nil, status.UpdateStatus(p.connErr)
	}

	// 认证
	if p.Etag != etag {
		desc := fmt.Sprintf("etag is different %s-%s", p.Etag, etag)
		return nil, status.NewStatusDesc(scode.PlayerAuthWrong, desc)
	}

	// 判断是否有消息推送,没有说明游戏没开始或者上次请求还未处理
	if len(p.msgChan) == 0 {
		desc := fmt.Sprintf("maybe game not start or last req is not finish")
		return nil, status.NewStatusDesc(scode.MessageChannelEmpty, desc)
	}

	// 获取本次消息，并异步等待响应
	msg := <-p.msgChan
	go func() {
		msg.pMsg = <-p.resp
		msg.Done()
	}()
	p.req = msg.gMsg
	return msg.gMsg, msg.err
}

// 回复game端信息
func (p *Player) PutCards(req *ipm.PutMessageRequest) error {
	// 收发信息出错
	if p.connErr != nil {
		return status.UpdateStatus(p.connErr)
	}

	if p.Etag != req.Etag {
		desc := fmt.Sprintf("etag is different %s-%s", p.Etag, req.Etag)
		return status.NewStatusDesc(scode.PlayerAuthWrong, desc)
	}

	// 检查出牌是否有效
	if !p.ownerCards(req.Pmsg.PutCards) || !p.checkRespValid(req.Pmsg) {
		desc := fmt.Sprintf("please check message type and card type value length")
		return status.NewStatusDesc(scode.MessageInvalid, desc)
	}

	p.resp <- req.Pmsg
	return nil
}

// 检查消息类型，值是否有效
func (p *Player) checkRespValid(resp *igm.PlayerMessage) bool {
	if p.req.RoundOwner == p.Id {
		// 争夺地主中只能选择抢或者过
		// 出牌判断是否可出
		switch p.req.MsgType {
		case igm.Over:
			return true
		case igm.Get:
			resp.PutCards = nil
			if resp.MsgType == igm.Pass {
				return true
			}
			if resp.MsgType != igm.Get {
				return false
			}
		case igm.Put:
			// 地主第一次不能不出牌,上次牌是自己的不能不出
			if resp.MsgType == igm.Pass && (p.req.LastId != "" || p.req.LastId != p.Id) {
				resp.PutCards = nil
				return true
			}
			if resp.MsgType != igm.Put {
				return false
			}

			repeated, value, length := card.GetRepeatNumAndValue(resp.PutCards)
			curType := card.GetCardsType(repeated, resp.PutCards)

			// 牌类型正确且上次牌无人要，或者地主首次出牌
			if curType != card.ValueTypeUnknown && (p.req.LastId == p.Id || p.req.LastId == "") {
				return true
			}

			lastRepeated, lastValue, lastLength := card.GetRepeatNumAndValue(p.req.LastCards)
			lastType := card.GetCardsType(lastRepeated, p.req.LastCards)

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
		resp.MsgType = igm.Over
		resp.PutCards = nil
	}

	return true
}

// 检查是否拥有这些牌
func (p *Player) ownerCards(cards []int64) bool {
	if len(cards) == 0 {
		return true
	}
	cMap := make(map[int64]struct{})
	for _, v := range p.req.YourCards {
		cMap[v] = struct{}{}
	}
	for _, v := range cards {
		if _, exist := cMap[v]; !exist {
			return false
		}
	}
	return true
}

// 建立数据连接
func (p *Player) ConnGame(etag, addr string) error {
	log.Infof("conn game:%s", addr)

	if p.Etag != etag {
		desc := fmt.Sprintf("etag is different")
		return status.NewStatusDesc(scode.PlayerAuthWrong, desc)
	}
	c, err := iclient.NewGameClient(addr)
	if err != nil {
		return status.UpdateStatus(err)
	}
	stream, err := c.PlayerConn(context.Background())
	if err != nil {
		desc := fmt.Sprintf("GRPC error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	pMsg := &igm.PlayerMessage{Pid: p.Id}
	err = stream.Send(pMsg)
	if err != nil {
		desc := fmt.Sprintf("GRPC Send error:%s", err)
		return status.NewStatusDesc(scode.GRPCError, desc)
	}
	p.gClient = c
	p.stream = stream
	p.connErr = nil
	go p.goListenGame()
	return nil
}

// 断开数据连接
func (p *Player) CloseConnGame(etag string) error {
	log.Infof("close game")

	if p.Etag != etag {
		desc := fmt.Sprintf("etag is different %s-%s", p.Etag, etag)
		return status.NewStatusDesc(scode.PlayerAuthWrong, desc)
	}
	if p.gClient == nil {
		desc := fmt.Sprintf("conn is nil")
		return status.NewStatusDesc(scode.PlayerNotAttachGame, desc)
	}
	close(p.close)
	if err := p.gClient.Close(); err != nil {
		log.Errorf("GRPC Close error:%s", err)
	}
	return nil
}

// 收发数据
func (p *Player) goListenGame() {
	for {
		msg := &message{}
		req, err := p.stream.Recv()
		if err != nil {
			desc := fmt.Sprintf("GRPC Recv error:%s", err)
			p.connErr = status.NewStatusDesc(scode.GRPCError, desc)
			return
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
			desc := fmt.Sprintf("GRPC Send error:%s", err)
			p.connErr = status.NewStatusDesc(scode.GRPCError, desc)
			return
		}
	}
}
