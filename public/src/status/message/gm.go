package message

import "scode"

var gmMsg = map[int32][2]string{
	scode.GameAlreadyExist:    {"game already exist", "该游戏已存在"},
	scode.GameNotExist:        {"game not exist", "该游戏不存在"},
	scode.GMDBOperateError:    {"GM DB operate error", "GM数据库操作错误"},
	scode.GMCallGoLibError:    {"GM call go lib error", "GM调用标准库函数报错"},
	scode.GRPCError:           {"GRPC error", "GRPC错误"},
	scode.GamePlayerIsNotHost: {"the player is not the game host", "该玩家不是房主"},
	scode.GameAlreadyStart:    {"game already start", "游戏已开始"},
	scode.GameAlreadyStop:     {"game already stop", "游戏已结束"},
	scode.NobodyWantGetLord:   {"nobody want to be the landlord", "无人想当地主"},
	scode.GamePlayerIsFull:    {"game room is full", "房间已满"},
	scode.PlayerNotInTheGame:  {"the player is not belong this game", "该玩家不属于该房间"},
	scode.GamePlayerNotEnough: {"game start need more player", "需要更多玩家加入才能开始游戏"},
}

func init() {
	registerMessage(gmMsg)
}
