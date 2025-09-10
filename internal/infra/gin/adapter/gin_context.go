package adapter

import (
	http "scheduling/internal/infra/gin"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	ctx *gin.Context
}

func (g *GinContext) Param(name string) string {
	return g.ctx.Param(name)
}
func (g *GinContext) Query(name string) string {
	return g.ctx.Query(name)
}
func (g *GinContext) Bind(obj interface{}) error {
	return g.ctx.ShouldBindJSON(obj)
}
func (g *GinContext) JSON(status int, obj interface{}) error {
	g.ctx.JSON(status, obj)
	return nil
}
func (g *GinContext) Status(code int) {
	g.ctx.Status(code)
}
func (g *GinContext) Header(key, value string) {
	g.ctx.Header(key, value)
}

func wrapHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = h(&GinContext{ctx: c})
	}
}

func wrapMiddleware(m http.MiddlewareFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = m(func(ctx http.Context) error {
			c.Next()
			return nil
		})(&GinContext{ctx: c})
	}
}
