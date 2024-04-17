package web

import "strings"

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
func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	// 判断是否已经存在,找到树
	root, ok := r.trees[method]
	if !ok {
		// 不存在,创建树
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	path = path[1:]
	// 切割 path
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 递归下去，找准位置
		// 中途有节点不存在需要创造出来
		children := root.childOrCreate(seg)
		root = children
	}
	root.handler = handleFunc
}

func (n *node) childOrCreate(seg string) *node {
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

//	type tree struct {
//		root *node
//	}
type node struct {
	path string
	// 子 path 到子节点的映射
	children map[string]*node
	handler  HandleFunc

	// 用户注册业务逻辑
}
