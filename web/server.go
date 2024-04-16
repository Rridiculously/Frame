package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保一定实现Server接口
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// AddRoute 路由注册功能
	// method 是 HTTP 方法
	// path 是路由
	// handleFunc 是处理函数业务逻辑
	AddRoute(method string, path string, handleFunc HandleFunc)
	// AddRoute1 实现多个
	//AddRoute1(method string, path string, handlers ...HandleFunc)
}

type HTTPServer struct {
}

// ServeHTTP 处理请求入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 查找路由， 应且执行命中业务逻辑
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {

}
func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {

}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}
func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

// Start func (h *HTTPServer) AddRoute1(method string, path string, handlers ...HandleFunc) {
//
// }
func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 让用户注册 after start 回调
	//执行业务所需前置条件

	return http.Serve(l, h)
}
