package guu

import (
	"log"
	"time"
)

// Logger 全局中间件，用来记录或打印请求处理的时间
func Logger() HandlerFunc{
	return func(c *Context){
		//请求开始时间
		t := time.Now()
		//等待请求处理
		c.Next()
		//处理结束
		log.Printf("[%d] %s cost %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}