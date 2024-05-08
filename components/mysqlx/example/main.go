package main

import (
	"context"
	"github.com/TremblingV5/box/components/internal/example"
	"github.com/TremblingV5/box/components/mysqlx"
	"log"
)

func main() {
	example.InitComponentConfig(
		"./components/mysqlx/example/config",
		"mysql",
		"mysql",
		mysqlx.Init,
	)

	db := mysqlx.GetDBClient(context.Background(), "default")
	log.Println(db)
}
