package ff

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix     string
	middleware []HandlerFunc
	parent     *RouterGroup
	engine     *Engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.route.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middleware = append(g.middleware, middlewares...)
}

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	g.GET(urlPattern, handler)
}

type Engine struct {
	*RouterGroup
	route         *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func New() *Engine {
	engine := &Engine{route: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
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
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = e
	e.route.handle(c)
}
