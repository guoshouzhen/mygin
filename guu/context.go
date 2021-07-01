package guu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H map的别名，方便使用
type H map[string]interface{}

// Context http请求响应上下文
type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path   string
	Method string
	Params map[string]string //动态路由参数
	//response info
	StatusCode int
	//middleware, contains all routergroup's
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1, //记录当前执行到第几个中间件
	}
}

// Next 表示等待执行其他的中间件或用户定义的handler，调用后转移控制权至下一个中间件或者用户定义的handler
func (c *Context) Next() {
	c.index++
	size := len(c.handlers)
	for ; c.index < size; c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm 表单中查询一个参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query URL中查询一个参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Status 设置http响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置http响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 设置文本形式的响应体
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetHeader("Encoding", "utf-8")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 设置json格式数据响应体
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 返回二进制流
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 返回html文本格式数据
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// Fail 请求失败，停止执行中间件及用户定义handler
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{
		"message": err,
	})
}
