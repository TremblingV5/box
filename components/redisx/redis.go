package redisx

import (
	"context"
	"fmt"
	"github.com/TremblingV5/box/components"
	"github.com/go-redis/redis"
	"sync"
)

var (
	globalClientMap = sync.Map{}
	globalConfigMap = make(components.ConfigMap[*Config])
)

type RedisClients struct {
	clients sync.Map
}

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) error {
	globalConfigMap = cm

	for k, v := range cm {
		db, err := Connect(v)

		if err != nil {
			return err
		}

		globalClientMap.Store(k, db)
	}

	return nil
}

func Connect(c *Config) (*RedisClients, error) {
	c.SetDefault()

	option := &redis.Options{}
	option.Addr = c.ToDSN()
	if c.Password != "" {
		option.Password = c.Password
	}

	clients := &RedisClients{}

	for _, item := range c.DBList {
		o := &redis.Options{
			Addr:     option.Addr,
			Password: option.Password,
			DB:       item.Number,
		}
		client := redis.NewClient(o)
		if _, err := client.Ping().Result(); err != nil {
			panic(err)
		}

		clients.clients.Store(item.Name, client)
	}

	return clients, nil
}

func GetRedisClient(ctx context.Context, storeKey, dbKey string) *redis.Client {
	if v, ok := globalClientMap.Load(storeKey); ok {
		if c, ok := v.(*RedisClients); ok {
			if v, ok := c.clients.Load(dbKey); ok {
				if client, ok := v.(*redis.Client); ok {
					return client
				}
			}
		}
	}

	panic(fmt.Sprintf("%s, %s not init", storeKey, dbKey))
}
