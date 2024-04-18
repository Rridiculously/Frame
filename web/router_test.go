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
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}
	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 验证路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": &node{
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "web: path already exists")
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
func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodDelete,
			path:   "/",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/user",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/user/home",
		//},
		{
			method: http.MethodDelete,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}
	r := newRouter()
	var mockHandler HandleFunc = func(ctx Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}
	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNNode *node
	}{
		{ // 完全命中
			name:      "order detail",
			method:    http.MethodDelete,
			path:      "/order/detail",
			wantFound: true,
			wantNNode: &node{
				path:    "detail",
				handler: mockHandler,
			},
		},
		{ // 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantNNode: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"order": &node{
						path: "order",

						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.wantNNode.path, n.path)
			//assert.Equal(t, tc.wantNNode.children, n.children)
			nHandler := reflect.ValueOf(n.handler)
			tHandler := reflect.ValueOf(tc.wantNNode.handler)
			assert.True(t, nHandler == tHandler)
		})
	}
}
