package gei

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const LogPrefix = "GEI:"

type HandlerFunc func(ctx *Context)

type RouterGroup struct {
	prefix     string        //分组前缀
	middleware []HandlerFunc //中间件
	parent     *RouterGroup
	engine     *Engine //RouterGroup也要有engine的功能
}
type Engine struct {
	*RouterGroup //嵌套
	router       *router
	groups       []*RouterGroup //储存所有组
}

// 实现Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

// New 创建一个新engine
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

// Group 创建一个分组路由
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine //全部组共用一个engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 创建一个route
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run 监听启动
func (engine *Engine) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}
	fmt.Println("Now Server Listening On:", addr)
	err := http.ListenAndServe(addr, engine)
	if err != nil {
		return errors.New("failed to start server")
	}
	return nil
}
