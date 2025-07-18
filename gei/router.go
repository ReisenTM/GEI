package gei

import (
	"fmt"
	"net/http"
	"time"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	router.handlers[method+"-"+pattern] = handler
}
func (router *router) handle(c *Context) {
	fmt.Printf("%s\t%s[%s] - %s\n", time.Now().Format("2006-01-02 15-01-05"), LogPrefix, c.Method, c.Path)
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
