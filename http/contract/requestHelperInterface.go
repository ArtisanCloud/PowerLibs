package contract

import (
	"net/http"
)

type RequestHandle func(request *http.Request) (response *http.Response, err error)

type RequestMiddleware func(handle RequestHandle) RequestHandle

// RequestHelperInterface 是一个 Client Wrap 类型接口, 面向上层
type RequestHelperInterface interface {
	SetClient(client ClientInterface)
	GetClient() ClientInterface

	// WithMiddleware 设置中间件, 在初始化 RequestDataflowInterface 时调用
	WithMiddleware(middlewares ...RequestMiddleware)

	// Df 返回链式构建实例, 该实例应该引用 RequestHelperInterface 的 ClientDriver 和 RequestHandle(middleware build以后的方法)
	Df() RequestDataflowInterface
}
