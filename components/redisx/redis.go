package redisx

import (
	"context"
	"fmt"
	"github.com/TremblingV5/box/components"

	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	globalClientMap = sync.Map{}
	globalConfigMap = make(components.ConfigMap[*Config])
)

func getGlobalClientMapKey(configKey, dbName string) string {
	return fmt.Sprintf("%s.%s", configKey, dbName)
}

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) error {
	globalConfigMap = cm

	for k, v := range cm {
		Connect(k, v)
	}

	return nil
}

func Connect(configKey string, c *Config) {
	c.SetDefault()

	option := &redis.Options{}
	option.Addr = c.ToDSN()
	if c.Password != "" {
		option.Password = c.Password
	}

	for name, number := range c.DBList {
		o := &redis.Options{
			Addr:     option.Addr,
			Password: option.Password,
			DB:       number,
		}
		client := redis.NewClient(o)
		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		globalClientMap.Store(getGlobalClientMapKey(configKey, name), client)
	}
}

// GetClient used to get a redis client instance
// keys is used to declare get which one
// Index 0 of keys is the store key
// Index 1 of keys is the db key
// If keys is empty, it will return the default client
func GetClient(ctx context.Context, keys ...string) *redis.Client {
	storeKey := "default"
	dbKey := "default"

	if len(keys) > 0 {
		storeKey = keys[0]
	}

	if len(keys) > 1 {
		dbKey = keys[1]
	}

	v, ok := globalClientMap.Load(getGlobalClientMapKey(storeKey, dbKey))
	if !ok {
		panic(fmt.Sprintf("%s, %s not init", storeKey, dbKey))
	}

	return v.(*redis.Client)
}
