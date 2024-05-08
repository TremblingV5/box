package hbasex

import (
	"context"
	"github.com/TremblingV5/box/components"
	"github.com/tsuna/gohbase"
	"sync"
)

var (
	globalClientMap = sync.Map{}
	globalConfigMap = make(components.ConfigMap[*Config])
)

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) error {
	globalConfigMap = cm

	for k, v := range cm {
		client, err := Connect(v)
		if err != nil {
			panic(err)
		}

		globalClientMap.Store(k, client)
	}

	return nil
}

func Connect(c *Config) (*gohbase.Client, error) {
	c.SetDefault()

	client := gohbase.NewClient(c.Host)
	return &client, nil
}

func getKey(keys ...string) string {
	if len(keys) == 0 {
		return "default"
	}

	return keys[0]
}

func GetClient(ctx context.Context, keys ...string) *gohbase.Client {
	key := getKey(keys...)
	client, ok := globalClientMap.Load(key)
	if !ok {
		panic("client not found")
	}

	return client.(*gohbase.Client)
}
