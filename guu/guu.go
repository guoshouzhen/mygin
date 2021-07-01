package guu

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc HandlerFunc defines the request handler used by guu, only for developer.
type HandlerFunc func(c *Context)

// RouterGroup 将所有跟路由有关的函数都交给RouterGroup实现
type RouterGroup struct {
	prefix      string        //前缀
	middlewares []HandlerFunc //支持中间件（路由分组的基础上应用）
	parent      *RouterGroup  //支持分组嵌套
	engine      *Engine       //all groups share a engine instance
}

// Engine Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup //使Engine拥有RouterGroup的所有能力
	//路由-函数映射表
	router *router
	groups []*RouterGroup // store all groups
}

// New Init Engine struct and return it's pointer
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 创建一个新的RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 将中间件应用到某个group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// private method, save router and it's handle func
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	//路由由其请求方法及路径拼接而成，可应对相同路径但不同请求方法的情况，例：POST-/hello
	group.engine.router.addRouter(method, pattern, handler)
}

// GET GET is method to add GET request handle func
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST POST is method to add POST request handle func
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// ServeHTTP Implement ServeHTTP method
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//中间件处理，整合前缀相同的所有分组路由下的中间件
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares //该请求需要应用的中间件
	engine.router.handle(c)
}

// Run app run, listen on the specify port.
func (engine *Engine) Run(addr string) (err error) {
	//addr为监听端口，第二个参数是处理所有http请求的实例（handler Handler），若为nil，则代表使用标准库中的实例处理
	//Handler是一个接口
	//type Handler interface {
	//	ServeHTTP(ResponseWriter, *Request)
	//}
	//使用一个实现了serverHTTP方法的实例engine，即可传入ListenAndServe并用该实例来处理所有请求
	return http.ListenAndServe(addr, engine)
}
