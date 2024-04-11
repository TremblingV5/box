package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/TremblingV5/box/httpserver/hertzx"
	"github.com/TremblingV5/box/launcher"
)

func registerHelloWorldHandler() hertzx.RegisterHertzRouter {
	return func(group *route.RouterGroup) {
		group.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(200, utils.H{
				"message": "hello world",
			})
		})
	}
}

func initHertzServer() *hertzx.HertzServer {
	hertzServer := hertzx.NewHertzServer(":8080", "/")
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
