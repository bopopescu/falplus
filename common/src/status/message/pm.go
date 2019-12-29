package message

import "scode"

var pmMsg = map[int32][2]string{
	scode.PMDBOperateError:    {"PM DB operate error","PM数据库操作错误"},
	scode.PlayerAlreadyExist:  {"player already exist","玩家已存在"},
	scode.PlayerNotExist:      {"player not exist","玩家不存在"},
	scode.PMCallGoLibError:    {"PM call go lib error","PM调用标准库出错"},
	scode.NameOrPasswordError: {"name or password error","账户名或密码错误"},
	scode.PlayerAuthWrong:     {"player authentication error","玩家认证错误"},
	scode.PlayerNotAttachGame: {"player not attach game","玩家尚未加入游戏"},
	scode.MessageChannelEmpty: {"message channel is empty","未获得新消息"},
	scode.MessageInvalid:      {"message is invalid","消息无效"},
	scode.PutMessageBeforeGet: {"put message before get","必须收到消息才能发送回复"},
}

func init(){
	registerMessage(pmMsg)
}
