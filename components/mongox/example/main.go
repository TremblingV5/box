package main

import (
	"context"
	"github.com/TremblingV5/box/components/internal/example"
	"github.com/TremblingV5/box/components/mongox"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func main() {
	example.InitComponentConfig(
		"./components/mongox/example/config",
		"mongo",
		"mongo",
		mongox.Init,
	)

	client := mongox.GetClient(context.Background())
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Println(err)
	}

	box1 := mongox.GetCollection(context.Background(), "default", "box1")
	log.Println(box1.Name())
}
