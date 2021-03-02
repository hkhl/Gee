package gee

import (
	"log"
	"net/http"
	"strings"
)

//type HandleFunc func(w http.ResponseWriter, r *http.Request)
type HandleFunc func(c *Context)

type RouterGroup struct {
	prefix	 	string
	middlewares []HandleFunc
	parent 		*RouterGroup
	engine 		*Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup  // store all groups
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//r := gee.New()
//v1 := r.Group("/v1")
func (group *RouterGroup)Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler  HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler  HandleFunc) {
	group.addRoute("GET", pattern, handler) //group.engine.addRouter = group.engine.RouterGroup.addRouter prefix为空
}

func (group *RouterGroup) POST(pattern string, handler  HandleFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(writer, request)
	c.handlers = middlewares
	engine.router.handle(c)
}
