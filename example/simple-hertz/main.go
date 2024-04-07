package main

import (
	"context"
	"github.com/TremblingV5/box/httpserver"
	"github.com/TremblingV5/box/launcher"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
)

func registerHelloWorldHandler() httpserver.RegisterHertzRouter {
	return func(group *route.RouterGroup) {
		group.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(200, utils.H{
				"message": "hello world",
			})
		})
	}
}

func initHertzServer() *httpserver.HertzServer {
	hertzServer := httpserver.NewHertzServer(":8080", "/")
	hertzServer.RegisterHttpHandlers(registerHelloWorldHandler())
	return hertzServer
}

func main() {
	l := launcher.New()

	l.AddBeforeServerStartHandler(func() {
		l.AddServer(initHertzServer())
	})

	l.Run()
}
