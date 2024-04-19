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
func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return root, true
	}
	// 去掉开头和结尾的 /
	path = strings.Trim(path, "/")

	// 切割 path
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	return root, true

}
func (n *node) childOrCreate(seg string) *node {
	if seg == "*" {
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
func (n *node) childOf(seg string) (*node, bool) {
	if n.children == nil {
		return n.starChild, n.starChild != nil
	}
	child, ok := n.children[seg]
	if !ok {
		return n.starChild, n.starChild != nil
	}
	return child, true
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
	// 用户注册业务逻辑
}
