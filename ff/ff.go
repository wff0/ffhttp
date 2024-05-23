package ff

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	route *router
}

func New() *Engine {
	return &Engine{
		route: newRoute(),
	}
}

func (e *Engine) addRoute(method, path string, handler HandlerFunc) {
	e.route.addRoute(method, path, handler)
}

func (e *Engine) GET(path string, handlerFunc HandlerFunc) {
	e.addRoute("GET", path, handlerFunc)
}

func (e *Engine) POST(path string, handlerFunc HandlerFunc) {
	e.addRoute("POST", path, handlerFunc)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.route.handle(c)
}
