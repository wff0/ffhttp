package ff

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'], roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc], handlers['POST-/p/:lang/doc']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0, len(vs))
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	parts := parsePattern(pattern)

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := fmt.Sprintf("%s-%s", method, pattern)
	r.handlers[key] = handler

}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	// 找到path对应的handler
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := fmt.Sprintf("%s-%s", c.Method, n.pattern)
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
