package mongox

import (
	"context"
	"fmt"
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
		globalClientMap.Store(k, Connect(k, v))
	}

	return nil
}

func Connect(configKey string, c *Config) *mongo.Client {
	opt := options.Client().ApplyURI(c.ToDSN())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}

	db := client.Database(c.Database)
	for _, collectionName := range c.Collections {
		globalCollectionMap.Store(fmt.Sprintf("%s_%s", configKey, collectionName), db.Collection(collectionName))
	}

	return client
}

func GetClient(ctx context.Context, keys ...string) *mongo.Client {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}

	if v, ok := globalClientMap.Load(key); ok {
		return v.(*mongo.Client)
	}

	return nil
}

// GetCollection used to get a collections by given keys
// Index 0 is the config key
// Index 1 is the collection name
// If not given, related parameters will be set to 'default'
func GetCollection(ctx context.Context, keys ...string) *mongo.Collection {
	configKey := "default"
	collectionName := "default"

	if len(keys) > 0 {
		configKey = keys[0]
	}

	if len(keys) > 1 {
		collectionName = keys[1]
	}

	if v, ok := globalCollectionMap.Load(fmt.Sprintf("%s_%s", configKey, collectionName)); ok {
		return v.(*mongo.Collection)
	}

	return nil
}
