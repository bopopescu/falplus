package scode

// GM错误码
const (
	GameAlreadyExist = 1001
	GameNotExist     = 1002
	GMDBOperateError = 1003
	GMCallGoLibError = 1004
	GRPCError        = 1005
)

// game错误码
const (
	GamePlayerIsNotHost = 2001
	GameAlreadyStart    = 2002
	GameAlreadyStop     = 2003
	NobodyWantGetLord   = 2004
	GamePlayerIsFull    = 2005
	PlayerNotInTheGame  = 2006
	GamePlayerNotEnough = 2007
)
