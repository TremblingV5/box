package mongox

import (
	"context"
	"github.com/TremblingV5/box/components"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	globalClientMap     = sync.Map{}
	globalCollectionMap = sync.Map{}
	globalConfigMap     = make(components.ConfigMap[*Config])
)

type MongoDBClient struct {
	client *mongo.Client
}

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) error {
	globalConfigMap = cm

	for k, v := range cm {
		v.SetDefault()
		globalClientMap.Store(k, Connect(v))
	}

	return nil
}

func Connect(c *Config) *mongo.Client {
	opt := options.Client().ApplyURI(c.ToDSN())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}

	db := client.Database(c.Database)
	for _, collectionName := range c.Collections {
		globalCollectionMap.Store(collectionName, db.Collection(collectionName))
	}

	return client
}
