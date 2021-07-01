package main

import (
	"guu"
	"log"
	"net/http"
	"time"
)

func main() {
	r := guu.New()

	/*===============分组路由测试==============*/
	////注册路由及对应处理方法
	//r.GET("/index", func(c *guu.Context){
	//	c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	//})
	//
	//v1 := r.Group("/v1")
	//{
	//	v1.GET("/", func(c *guu.Context){
	//		c.HTML(http.StatusOK, "<h1>This is guu</h1>")
	//	})
	//
	//	v1.GET("/hello", func(c *guu.Context){
	//		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"),c.Path)
	//	})
	//
	//	v11 := v1.Group("/p1")
	//	{
	//		v11.GET("/sayhello",func(c *guu.Context){
	//			c.String(http.StatusOK,"say: %s\n", "hello")
	//		})
	//	}
	//}
	//
	//v2 := r.Group("/v2")
	//{
	//	v2.GET("/hello/:name",func(c *guu.Context){
	//		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//	})
	//
	//	v2.POST("/login", func(c *guu.Context) {
	//		c.JSON(http.StatusOK, guu.H{
	//			"username": c.PostForm("username"),
	//			"password": c.PostForm("password"),
	//		})
	//	})
	//}


	/*===============中间件测试==============*/
	r.Use(guu.Logger()) //global middleware
	r.GET("/", func(c *guu.Context){
		c.HTML(http.StatusOK, "<h1>This is guu</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *guu.Context){
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	//监听端口
	r.Run(":9999")

}


//测试中间件
func onlyForV2() guu.HandlerFunc{
	return func(c *guu.Context){
		//start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
