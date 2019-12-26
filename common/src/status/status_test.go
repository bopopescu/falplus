package status

import (
	"errors"
	"fmt"
	"scode"
	"testing"
)

func TestNewError(t *testing.T) {
	testError(NewStatus(scode.ScodeComponentAlreadyExist, nil, "vid-1-1"))
	/*
		code:  8010
		message:  component vid-1-1 already exist
		messageCn:  组件vid-1-1已经存在
	*/

	testError(NewStatus(scode.ScodeComponentNotExist, nil, "vid-1-2"))
	/*
		code:  8011
		message:  component vid-1-2 do not exist
		messageCn:  组件vid-1-2不存在
	*/

	testError(NewStatus(scode.ScodeComponentTypeError, nil))
	/*
		code:  8030
		message:  component type error
		messageCn:  组件类型错误
	*/

	testError(NewStatus(scode.ScodeComponentCreateError, nil))
	/*
		code:  8031
		message:  Error message is not defined
		messageCn:  错误信息没有定义
	*/

	testError(errors.New("normal error"))
	/*
		code:  -1
		message:  invalid error message
		messageCn:  无效的错误信息
	*/

	testError(nil)
	/*
		code:  0
		message:
		messageCn:
	*/
}

func testError(err error) {
	fmt.Println("code: ", Code(err))
	fmt.Println("message: ", Message(err))
	fmt.Println("messageCn: ", MessageCn(err))
}

func TestToError(t *testing.T) {
	testToError(NewStatus(scode.ScodeComponentAlreadyExist, nil, "vid-1-1"))
	testToError(NewStatus(scode.ScodeComponentNotExist, nil, "vid-1-2"))
	testToError(NewStatus(scode.ScodeComponentTypeError, nil))
	testToError(NewStatus(scode.ScodeComponentCreateError, nil))
	testToError(errors.New("normal error"))
	testToError(nil)
}

func testToError(err error) {
	Err := FromError(err)
	fmt.Println("code: ", Err.Code)
	fmt.Println("message: ", Err.Message)
	fmt.Println("messageCn: ", Err.MessageCn)
}

func TestStack(t *testing.T) {
	err1 := NewStatus(scode.ScodeComponentAlreadyExist, nil, "vid-1-1")
	err2 := NewStatus(scode.ScodeComponentNotExist, err1, "vid-1-2")
	err3 := NewStatusStack(scode.ScodeComponentTypeError, err2, nil)
	testStack(UpdateStatus(FromError(err3)))
	fmt.Println(UpdateStatus(FromError(err3)).GetStackInfo())
	fmt.Println(UpdateStatus(FromError(err3)))
}

func testStack(err error) {
	Err := FromError(err)
	fmt.Println(Err.Stack)
}
