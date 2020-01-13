package scode

// PM错误码
const (
	PMDBOperateError    = 3001
	PlayerAlreadyExist  = 3002
	PlayerNotExist      = 3003
	PMCallGoLibError    = 3004
	NameOrPasswordError = 3005
	PlayerAuthWrong     = 3006
)

// player错误码
const (
	PlayerNotAttachGame = 4001
	MessageChannelEmpty = 4002
	MessageInvalid      = 4003
	PutMessageBeforeGet = 4004
)
