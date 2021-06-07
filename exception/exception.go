package exception

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
)

type Exception struct {
}

func (e *Exception) HandleException(ctx context.Context, action string, args ...interface{}) {
	if p := recover(); p != nil {
		var err error
		switch rs := p.(type) {

		case runtime.Error:
			err = p.(runtime.Error)
		case string:
			err = errors.New(p.(string))
		// 若非APIResponse，也许默认抛出一个若非APIResponse
		default:
			fmt.Printf("%v \n", rs)
		}

		if err != nil {
			fmt.Printf("panic: %v \r\n", err.Error())
			fmt.Printf("err stack: %v \r\n", string(debug.Stack()))
		}
	}
}

func (e *Exception) GetMessage() string {
	return ""
}

func (e *Exception) GetCode() int {
	return 0
}

func (e *Exception) GetFile() string {
	return ""
}

func (e *Exception) GetLine() int {
	return 0
}

func (e *Exception) GetTrace() []string {
	return nil
}

func (e *Exception) GetTraceAsString() string {
	return ""
}

func (e *Exception) GetPrevious() *Throwable {
	return nil
}

func (e *Exception) __toString() string {
	return ""
}
