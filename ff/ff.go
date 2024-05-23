package ff

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	route map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{
		route: make(map[string]HandlerFunc),
	}
}

func (e *Engine) addRoute(method, path string, handlerFunc HandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, path)
	e.route[key] = handlerFunc
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
	key := fmt.Sprintf("%s-%s", r.Method, r.URL.Path)
	if handler, ok := e.route[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}
