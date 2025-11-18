package http

type Context interface {
	Param(name string) string
	Query(name string) string
	Bind(obj any) error
	JSON(status int, obj any) error
	Status(code int)
	Header(key, value string)
	Set(key string, value any)
}

type Router interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	Use(middleware ...MiddlewareFunc)
	Run(addr string) error
	Group(path string) Router
}

type HandlerFunc func(Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc
