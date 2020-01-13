package status

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"status/message"
	"strconv"
	"strings"
	"util"
)

var (
	progressInfo  = fmt.Sprintf("host:%s,pid:%d,module:%s", util.GetIPv4Addr(), os.Getpid(), filepath.Base(os.Args[0]))
	SuccessStatus = &Status{
		Code: 0,
	}
)

type StackInfo struct {
	host     string
	pid      string
	module   string
	code     string
	filename string
	funname  string
	line     string
}

func NewStatus(sCode int32, a ...interface{}) *Status {
	msg := message.StatusMsg[sCode]
	stack := []string{progressInfo + fmt.Sprintf(",code:%d,", sCode) + Caller()}
	return &Status{
		Code:      sCode,
		Message:   fmt.Sprintf(msg[0], a...),
		MessageCn: fmt.Sprintf(msg[1], a...),
		Stack:     stack,
		UUID:      uuid.New().String(),
	}
}

func NewStatusDesc(sCode int32, desc string, a ...interface{}) *Status {
	msg := message.StatusMsg[sCode]
	stack := []string{progressInfo + fmt.Sprintf(",code:%d,", sCode) + Caller()}
	return &Status{
		Code:      sCode,
		Message:   fmt.Sprintf(msg[0], a...),
		MessageCn: fmt.Sprintf(msg[1], a...),
		Stack:     stack,
		Desc:      desc,
		UUID:      uuid.New().String(),
	}
}

func NewStatusStack(sCode int32, sStatus error, a ...interface{}) *Status {
	msg := message.StatusMsg[sCode]
	stack := []string{progressInfo + fmt.Sprintf(",code:%d,", sCode) + Caller()}
	desc := ""
	if sStatus != nil {
		if st, ok := sStatus.(*Status); ok {
			stack = append(stack, st.Stack...)
			desc = st.Desc
		}
	}
	return &Status{
		Code:      sCode,
		Message:   fmt.Sprintf(msg[0], a...),
		MessageCn: fmt.Sprintf(msg[1], a...),
		Stack:     stack,
		Desc:      desc,
		UUID:      uuid.New().String(),
	}
}

// Error is for the error interface
func (st Status) Error() string {
	siList := st.GetStackInfo()
	callChain := ""
	lastService := ""
	for i, si := range siList {
		service := fmt.Sprintf("([%s:%s]", si.host, si.module)
		if i == 0 {
			callChain = fmt.Sprintf("%s<%s:%s>)", si.funname, si.filename, si.line)
		} else if service != lastService && lastService != "" {
			callChain = si.funname + ")->" + lastService + callChain
		} else {
			callChain = si.funname + "->" + callChain
		}
		lastService = service
	}
	callChain = lastService + callChain
	if st.Desc != "" {
		return fmt.Sprintf("UUID:%s,stack:%s,code:%d,message:%s,desc:%s", st.UUID, callChain, st.Code, st.Message, st.Desc)
	} else {
		return fmt.Sprintf("UUID:%s,stack:%s,code:%d,message:%s", st.UUID, callChain, st.Code, st.Message)
	}
}

func (st Status) GetStackInfo() []StackInfo {
	siList := make([]StackInfo, 0)
	for _, str := range st.Stack {
		si := StackInfo{}
		slist := strings.Split(str, ",")
		for _, s := range slist {
			if strings.HasPrefix(s, "host:") {
				si.host = strings.TrimPrefix(s, "host:")
			}
			if strings.HasPrefix(s, "pid:") {
				si.pid = strings.TrimPrefix(s, "pid:")
			}
			if strings.HasPrefix(s, "module:") {
				si.module = strings.TrimPrefix(s, "module:")
			}
			if strings.HasPrefix(s, "code:") {
				si.code = strings.TrimPrefix(s, "code:")
			}
			if strings.HasPrefix(s, "filename:") {
				si.filename = strings.TrimPrefix(s, "filename:")
			}
			if strings.HasPrefix(s, "func:") {
				si.funname = strings.TrimPrefix(s, "func:")
			}
			if strings.HasPrefix(s, "line:") {
				si.line = strings.TrimPrefix(s, "line:")
			}
		}
		siList = append(siList, si)
	}
	return siList
}

func UpdateStatus(err error) *Status {
	st := FromError(err)
	if st != nil {
		st.Stack = append(st.Stack, progressInfo+fmt.Sprintf(",code:%d,", st.Code)+Caller())
	}
	return st
}

func Caller() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	funcName := "???"
	f := runtime.FuncForPC(pc)
	if f != nil {
		funcName = f.Name()
	}
	_, filename := path.Split(file)
	flist := strings.Split(funcName, ".")
	funcName = flist[len(flist)-1]
	format := "filename:" + filename + ",func:" + funcName + ",line:" + strconv.FormatInt(int64(line), 10)
	return format
}

func FromError(err error) *Status {
	if err == nil {
		return &Status{
			Code: 0,
		}
	}
	if st, ok := err.(*Status); ok {
		return st
	} else {
		stack := []string{progressInfo + fmt.Sprintf(",code:%d,", -1) + Caller()}
		return &Status{
			Code:      -1,
			Message:   "invalid status message",
			MessageCn: "无效的状态信息",
			Stack:     stack,
		}
	}
}

func Code(err error) int32 {
	if err == nil {
		return 0
	}
	if st, ok := err.(*Status); ok {
		return st.Code
	} else {
		return -1
	}
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	if st, ok := err.(*Status); ok {
		return st.Message
	} else {
		return "invalid status message"
	}
}

func MessageCn(err error) string {
	if err == nil {
		return ""
	}
	if st, ok := err.(*Status); ok {
		return st.MessageCn
	} else {
		return "无效的状态信息"
	}
}

func Stack(err error) []string {
	if err == nil {
		return nil
	}
	if st, ok := err.(*Status); ok {
		return st.Stack
	} else {
		return nil
	}
}
