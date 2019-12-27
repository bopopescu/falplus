package player

import (
	"api/igm"
	"api/ipm"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
)

type Player struct {
	Id      string
	User    string
	Pwd     string
	Score   int
	Cards   []int
	conn    *grpc.ClientConn
	stream  igm.Game_PlayerConnClient
	close   chan struct{}
	msgChan chan *message
	resp    chan *igm.PlayerMessage
	req     *igm.GameMessage
	etag    string
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
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Player) UpdateInfo(info *ipm.PlayerInfo) {
	info.Etag = uuid.New().String()
	p.etag = info.Etag
	p.User = info.Name
	p.Pwd = info.Password
	p.port = info.Port
	p.Id = info.Id
}

func (p *Player) ShowGameMsg(etag string) (*igm.GameMessage, error) {
	if p.etag != etag {
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
	if p.etag != req.Etag {
		return fmt.Errorf("auth error")
	}
	if !p.checkRespValid(req.Pmsg) {
		return fmt.Errorf("the resp is invalid")
	}
	p.resp <- req.Pmsg
	return nil
}

func (p *Player) checkRespValid(resp *igm.PlayerMessage) bool {
	return false
}

func (p *Player) ConnGame(etag, addr string) error {
	if p.etag != etag {
		return fmt.Errorf("auth error")
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return err
	}
	c := igm.NewGameClient(conn)
	stream, err := c.PlayerConn(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}
	p.conn = conn
	p.stream = stream
	go p.ListenGame()
	return nil
}

func (p *Player) CloseConnGame(etag string) error {
	if p.etag != etag {
		return fmt.Errorf("auth error")
	}
	if p.conn == nil {
		return fmt.Errorf("conn is nil")
	}
	return p.conn.Close()
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
