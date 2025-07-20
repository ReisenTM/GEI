package gei

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// parsePattern 解析pattern
func parsePattern(pattern string) []string {
	todoList := strings.Split(pattern, "/")
	var parts []string
	for _, part := range todoList {
		if part != "" {
			parts = append(parts, part)
			if strings.HasPrefix(part, "*") {
				break
			}
		}
	}
	return parts
}

// newRouter 总路由管理
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// addRoute 新增路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{} //创建该key对应的根节点
	}
	r.roots[method].insert(pattern, parts, 0) //插入节点(height在插入过程中完成自增)
	r.handlers[key] = handler                 //绑定handler函数
}

// getRoute 通过传入的path找到对应的节点
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	params := make(map[string]string)
	searchParts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern) //找到对应的pattern
		for index, part := range parts {
			if part[0] == ':' {
				//填充
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' {
				//拼接后续path
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// handle 事务处理函数
func (r *router) handle(c *Context) {
	fmt.Printf("%s\t%s[%s] - %s\n", time.Now().Format("2006-01-02 15-01-05"), LogPrefix, c.Method, c.Path)
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//执行处理函数
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
