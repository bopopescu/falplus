package igm

const (
	Over = 0 // 无牌权时回复消息类型
	Get  = 1 // 抢地主
	Put  = 2 // 出牌
	Pass = 3 // 过
	Gpass = 4// 流局
)

var MsgType = map[int64]string{
	0:"over",
	1:"get",
	2:"put",
	3:"pass",
}

// 游戏状态 进程状态
const (
	Stop = 0
	Start = 1
)

type Game interface {
	AddPlayer(pid string) error
	DelPlayer(pid string) error
	PlayerConn(stream Game_PlayerConnServer) (<-chan struct{}, error)
	Start(pid string) error
	Stop(pid string) error
	State() int
}
