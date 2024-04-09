package main

import (
	"github.com/TremblingV5/box/configx"
	"github.com/TremblingV5/box/httpserver"
	"github.com/TremblingV5/box/launcher"
	"github.com/gin-gonic/gin"
)

func registerHelloWorldHandler() httpserver.RegisterGinRouter {
	return func(group *gin.RouterGroup) {
		group.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
	}
}

func initGinServer() *httpserver.GinServer {
	ginServer := httpserver.NewGinServer(":8080", "debug", "/")
	ginServer.RegisterHttpHandlers(registerHelloWorldHandler())
	return ginServer
}

func main() {
	l := launcher.New()

	l.AddBeforeConfigInitHandler(func() {
		configx.SetRootConfigPath("./components/mysqlx/example/config")
	})

	l.AddBeforeServerStartHandler(func() {
		l.AddServer(initGinServer())
	})

	l.Run()
}
