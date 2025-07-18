package gei

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const LogPrefix = "GEI:"

type HandlerFunc func(ctx *Context)
type Engine struct {
	router *router
}

// 增加路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// GET GET请求封装
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST POST请求封装
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 实现Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

// New 创建一个新engine
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
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
