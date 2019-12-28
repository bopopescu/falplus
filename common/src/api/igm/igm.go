package igm

const (
	Over = 0 // 无牌权时回复消息类型
	Get  = 1 // 抢地主
	Put  = 2 // 出牌
	Pass = 3 // 过
)

var MsgType = map[int64]string{
	0:"over",
	1:"get",
	2:"put",
	3:"pass",
}