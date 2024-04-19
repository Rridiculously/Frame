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
		//{
		//	method: http.MethodPost,
		//	path:   "/order/create",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/order/*",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/login",
		//},
	}
	var mockHandler HandleFunc = func(ctx *Context) {}
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
			//http.MethodPost: &node{
			//	path: "/",
			//	children: map[string]*node{
			//		"order": &node{
			//			path: "order",
			//			children: map[string]*node{
			//				"detail": &node{
			//					path:    "detail",
			//					handler: mockHandler,
			//				},
			//			},
			//			starChild: &node{
			//				path:    "*",
			//				handler: mockHandler,
			//			},
			//		},
			//		"login": &node{
			//			path:    "login",
			//			handler: mockHandler,
			//		},
			//	},
			//},
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

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}
	}
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
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
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}
	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		info      *matchInfo
	}{
		{ // 完全命中
			name:      "order detail",
			method:    http.MethodDelete,
			path:      "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "detail",
				},
			},
		},
		{ // 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "detail",
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
			assert.Equal(t, tc.info.n.path, n.n.path)
			//assert.Equal(t, tc.wantNNode.children, n.children)
			nHandler := reflect.ValueOf(n.n.handler)
			tHandler := reflect.ValueOf(tc.info.n.handler)
			assert.True(t, nHandler == tHandler)
		})
	}
}

//func TestRouter_FindRoute(t *testing.T) {
//	// 初始化你的路由器
//	r := newRouter()
//
//	// 定义路由
//	var mockHandler HandleFunc = func(ctx *Context) {}
//	r.addRoute("GET", "/users/:id", mockHandler)
//	r.addRoute("POST", "/posts/:id/comments/:commentID", mockHandler)
//	r.addRoute("GET", "/about", mockHandler)
//	// 根据需要添加更多路由
//
//	// 测试用例
//	tests := []struct {
//		method   string
//		path     string
//		expected *matchInfo
//		ok       bool
//	}{
//		// 测试案例1：有效路由
//		{
//			method: "GET",
//			path:   "/users/123",
//			expected: &matchInfo{
//				n: &node{
//					path: ":id",
//					// 添加预期的节点属性
//				},
//				pathParams: map[string]string{
//					"id": "123",
//				},
//			},
//			ok: true,
//		},
//		// 测试案例2：带多个参数的有效路由
//		{
//			method: "POST",
//			path:   "/posts/456/comments/789",
//			expected: &matchInfo{
//				n: &node{
//					path: "/posts/:id/comments/:commentID",
//					// 添加预期的节点属性
//				},
//				pathParams: map[string]string{
//					"id":        "456",
//					"commentID": "789",
//				},
//			},
//			ok: true,
//		},
//		// 测试案例3：路由未找到
//		{
//			method:   "GET",
//			path:     "/products",
//			expected: nil,
//			ok:       false,
//		},
//		// 根据需要添加更多测试用例
//	}
//
//	// 运行测试用例
//	for _, test := range tests {
//		matchInfo, ok := r.findRoute(test.method, test.path)
//		if ok != test.ok {
//			t.Errorf("期望方法 %s 和路径 %s 的 OK 为 %v，但得到 %v", test.method, test.path, test.ok, ok)
//		}
//		if ok && (matchInfo.n.path != test.expected.n.path || !comparePathParams(matchInfo.pathParams, test.expected.pathParams)) {
//			t.Errorf("期望方法 %s 和路径 %s 的匹配信息为 %v，但得到 %v", test.method, test.path, test.expected, matchInfo)
//		}
//	}
//}
//
//func comparePathParams(params1, params2 map[string]string) bool {
//	if len(params1) != len(params2) {
//		return false
//	}
//	for key, value := range params1 {
//		if params2Value, ok := params2[key]; !ok || params2Value != value {
//			return false
//		}
//	}
//	return true
//}
