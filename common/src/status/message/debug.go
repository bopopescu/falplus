package message

import "scode"

var debugMsg = map[int32][2]string{
	scode.DebugCallGoLibError:  {"debug call go lib error", "debug调用标准库出错"},
	scode.DebugAlreadyStop:		{"debug already stop", "debug已关闭"},
}

func init(){
	registerMessage(debugMsg)
}
