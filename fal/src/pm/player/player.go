package player

import (
	"api/igm"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
)

type Player struct {
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
}

var player *Player

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

func (p *Player) ShowCards() (*igm.GameMessage, error) {
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

func (p *Player) PutCards(resp *igm.PlayerMessage) error {
	if !p.checkRespValid(resp) {
		return fmt.Errorf("the resp is invalid")
	}
	p.resp <- resp
	return nil
}

func (p *Player) checkRespValid(resp *igm.PlayerMessage) bool {
	return false
}

func (p *Player) ConnGame(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := igm.NewGameClient(conn)
	stream, err := c.PlayerConn(context.Background())
	if err != nil {
		return err
	}
	p.conn = conn
	p.stream = stream
	return nil
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
