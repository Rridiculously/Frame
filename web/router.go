package web

import (
	"fmt"
	"strings"
)

// 支持对路由树操作
// 森林
type router struct {
	// 每个HTTPMethod对应一棵树
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

// addRoute 限制 Path 必须以 / 开头， 不能 / 结尾，。也不能有连续的 / 例如： ///
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" || path[0] != '/' {
		panic("web: path must start with '/'")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: path must not end with '/'")
	}
	if handleFunc == nil {
		panic("web: handleFunc is nil")
	}
	if method == "" {
		panic("web: method is empty")
	}

	// 判断是否已经存在,找到树
	root, ok := r.trees[method]
	if !ok {
		// 不存在,创建树
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 根节点特殊处理
	if path == "/" {
		if root.handler != nil {
			panic(fmt.Sprintf("web: path '%s' already exists", path))
		}
		root.handler = handleFunc
		return
	}

	// 切割 path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		// 校验连续的 /
		if seg == "" {
			panic("web: path segment is empty")
		}
		// 递归下去，找准位置
		// 中途有节点不存在需要创造出来
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: path '%s' already exists", path))
	}
	root.handler = handleFunc
}
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	// 去掉开头和结尾的 /
	path = strings.Trim(path, "/")

	// 切割 path
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 命中路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}

	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true

}
func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	res, ok := n.children[seg]
	if !ok {
		// 新建
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

// childOf 优先考虑静态匹配，匹配不上在考虑通配符匹配
// 返回值1：子节点 2：是否路径参数 3： 命中没有
func (n *node) childOf(seg string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[seg]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, true
}

//	type tree struct {
//		root *node
//	}
type node struct {
	path string

	// 静态匹配节点
	// 子 path 到子节点的映射
	children map[string]*node
	handler  HandleFunc

	// 通配符匹配
	starChild *node

	// 路径参数
	paramChild *node
	// 用户注册业务逻辑
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
