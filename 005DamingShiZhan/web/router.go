package web

// 路由树 (森林
type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

type node struct {
	path string

	// path 到子节点的映射
	children map[string]*node

	// 代表用户注册的业务逻辑
	handler HandleFunc
}

func (r *router) AddRoute(
	method string, path string, handler HandleFunc) {

}
