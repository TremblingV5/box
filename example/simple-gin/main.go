package main

import (
	"github.com/gin-gonic/gin"

	"github.com/TremblingV5/box/httpserver/ginx"
	"github.com/TremblingV5/box/launcher"
)

func registerHelloWorldHandler() ginx.RegisterGinRouter {
	return func(group *gin.RouterGroup) {
		group.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
	}
}

func initGinServer() *ginx.GinServer {
	ginServer := ginx.NewGinServer(":8080", "debug", "/")
	ginServer.RegisterHttpHandlers(registerHelloWorldHandler())
	return ginServer
}

func main() {
	l := launcher.New()

	l.AddBeforeServerStartHandler(func() {
		l.AddServer(initGinServer())
	})

	l.Run()
}
