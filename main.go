package main

import (
	"geekgo-webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	// server 处理GET请求 GET is a shortcut for router.Handle("GET", path, handlers)
	//server.GET(relativePath, HandlerFunc) 对于给定的请求路径，以HandlerFunc来处理请求返回响应
	// 示例
	//  r.GET("/ping", func(c *gin.Context) {
	//    c.JSON(http.StatusOK, gin.H{
	//      "message": "pong",
	//    })
	//  })
	u.RegisterRoutes(server)
	server.Run(":8080") // 网络编程

}
