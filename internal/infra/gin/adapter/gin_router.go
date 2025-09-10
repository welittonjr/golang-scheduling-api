package adapter

import (
	"fmt"

	http "scheduling/internal/infra/gin"

	"github.com/gin-gonic/gin"
)

type ginRouter struct {
	group gin.IRoutes
	root  *gin.Engine
}

func NewRouter() http.Router {
	engine := gin.Default()
	return &ginRouter{
		group: engine,
		root:  engine,
	}
}

func (r *ginRouter) GET(path string, handler http.HandlerFunc) {
	r.group.GET(path, wrapHandler(handler))
}

func (r *ginRouter) POST(path string, handler http.HandlerFunc) {
	r.group.POST(path, wrapHandler(handler))
}

func (r *ginRouter) PUT(path string, handler http.HandlerFunc) {
	r.group.PUT(path, wrapHandler(handler))
}

func (r *ginRouter) DELETE(path string, handler http.HandlerFunc) {
	r.group.DELETE(path, wrapHandler(handler))
}

func (r *ginRouter) Use(middlewares ...http.MiddlewareFunc) {
	for _, m := range middlewares {
		r.group.Use(wrapMiddleware(m))
	}
}

func (r *ginRouter) Group(path string) http.Router {
	if group, ok := r.group.(*gin.Engine); ok {
		newGroup := group.Group(path)
		return &ginRouter{group: newGroup, root: r.root}
	}
	if group, ok := r.group.(*gin.RouterGroup); ok {
		newGroup := group.Group(path)
		return &ginRouter{group: newGroup, root: r.root}
	}
	panic("unsupported gin.IRoutes type in Group()")
}

func (r *ginRouter) Run(addr string) error {
	if r.root != nil {
		return r.root.Run(addr)
	}
	return fmt.Errorf("no root engine to run")
}
