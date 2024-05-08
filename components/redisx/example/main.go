package main

import (
	"context"
	"github.com/TremblingV5/box/components/internal/example"
	"github.com/TremblingV5/box/components/mysqlx"
	"github.com/TremblingV5/box/components/redisx"
)

func main() {
	example.InitComponentConfig(
		"./components/redisx/example/config",
		"redis",
		"redis",
		mysqlx.Init,
	)

	client := redisx.GetRedisClient(context.Background(), "default", "default")
	if err := client.Ping().Err(); err != nil {
		panic(err)
	}
}
