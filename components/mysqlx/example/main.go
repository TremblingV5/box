package main

import (
	"github.com/TremblingV5/box/components"
	"github.com/TremblingV5/box/components/mysqlx"
	"github.com/TremblingV5/box/configx"
	"github.com/TremblingV5/box/launcher"
)

func main() {
	configx.SetRootConfigPath("./components/mysqlx/example/config")
	configx.ResolveWatcher(
		configx.Watch("database", configx.WithModel(&configx.Config{})),
	)
	configx.Init()

	store := configx.GetStore(configx.StoreKeyConfig)
	loadMap := make(configx.ComponentLoadMap)
	store.UnmarshalKey("componentLoad", &loadMap)

	l := launcher.NewComponentsLauncher(loadMap)
	l.CustomLaunch(func(l *launcher.ComponentsLauncher) {
		l.LaunchComponent("mysql", "mysql", func(storeKey, configKey string) error {
			return components.Load(storeKey, configKey, mysqlx.Init).Start()
		})
	})
}
