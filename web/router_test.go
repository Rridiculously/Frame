package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	// 构造路由树
	// 验证路由书
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "/",
		//},
	}
	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	// 验证路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintln("找不到对应 HTTPMethod"), false
		}
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}

	}
	return "", true
}
func (n *node) equal(y *node) (string, bool) {
	if y.path != n.path {
		return fmt.Sprintln("path 不一致"), false
	}

	if len(y.children) != len(n.children) {
		return fmt.Sprintln("子节点数量不一致"), false
	}
	// 比较 handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintln("handler 不一致"), false
	}
	for path, child := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintln("找不到对应子节点"), false
		}
		msg, equal := child.equal(dst)
		if !equal {
			return fmt.Sprintln(msg), false
		}
	}
	return "", true
}
