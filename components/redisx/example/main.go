package main

import (
	"context"
	"github.com/TremblingV5/box/components"
	"github.com/TremblingV5/box/components/redisx"
	"github.com/TremblingV5/box/configx"
	"github.com/TremblingV5/box/launcher"
)

func main() {
	configx.SetRootConfigPath("./components/redisx/example/config")
	configx.ResolveWatcher(
		configx.Watch("database", configx.WithModel(&configx.Config{})),
	)
	configx.Init()

	store := configx.GetStore(configx.StoreKeyConfig)
	loadMap := make(configx.ComponentLoadMap)
	store.UnmarshalKey("componentLoad", &loadMap)

	l := launcher.NewComponentsLauncher(loadMap)
	l.CustomLaunch(func(l *launcher.ComponentsLauncher) {
		l.LaunchComponent("redis", "redis", func(storeKey, configKey string) error {
			return components.Load(storeKey, configKey, redisx.Init).Start()
		})
	})

	client := redisx.GetRedisClient(context.Background(), "default", "default")
	if err := client.Ping().Err(); err != nil {
		panic(err)
	}
}
