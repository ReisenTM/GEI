package gei

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// context相关很复杂，尽量隐藏在内部，暴露的接口简单化

// H 是一个别名
type H map[string]interface{}
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	//Request
	Path   string
	Method string
	//Params
	//test/:lang/show->test/go/show   ':': ["lang"]:"go"
	//test/*lang/show->test/go/show   '*': ["lang"]:"go/show"
	Params map[string]string
	//Response
	StatusCode int
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *Context) Header(key string) string {
	return c.Writer.Header().Get(key)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (c *Context) XML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/xml")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (c *Context) Data(code int, contentType string, data []byte) {
	c.SetHeader("Content-Type", contentType)
	c.Status(code)
	c.Writer.Write(data)
}
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
