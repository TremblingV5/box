package launcher

import (
	"context"
	"github.com/TremblingV5/box/components"
	"github.com/TremblingV5/box/components/mysqlx"
	"github.com/TremblingV5/box/components/redisx"
	"github.com/TremblingV5/box/configx"
	"github.com/TremblingV5/box/gofer"
	"github.com/TremblingV5/box/logx"
	"go.uber.org/zap"
)

type ComponentsLauncher struct {
	group             *gofer.Group
	components        map[string]func() error
	componentsLoadMap configx.ComponentLoadMap
}

func NewComponentsLauncher(componentsLoadMap configx.ComponentLoadMap) *ComponentsLauncher {
	return &ComponentsLauncher{
		group: gofer.NewGroup(
			context.Background(),
			gofer.UseErrorGroup(),
		),
		components:        make(map[string]func() error),
		componentsLoadMap: componentsLoadMap,
	}
}

func (l *ComponentsLauncher) launch() {
	l.LaunchComponent("mysql", "mysql", func(storeKey, configKey string) error {
		return components.Load(storeKey, configKey, mysqlx.Init).Start()
	})

	l.LaunchComponent("redis", "redis", func(storeKey, configKey string) error {
		return components.Load(storeKey, configKey, redisx.Init).Start()
	})

	for _, fn := range l.components {
		l.group.Run(fn)
	}

	if err := l.group.Wait(); err != nil {
		panic("launch components error: " + err.Error())
	}
}

func (l *ComponentsLauncher) CustomLaunch(launch func(l *ComponentsLauncher)) {
	launch(l)

	for _, fn := range l.components {
		l.group.Run(fn)
	}

	if err := l.group.Wait(); err != nil {
		panic("launch components error: " + err.Error())
	}
}

func (l *ComponentsLauncher) LaunchComponent(name, configKey string, launch func(storeKey, configKey string) error) {
	storeKey := configx.StoreKeyDefault

	if loadConfig, ok := l.componentsLoadMap[name]; ok && loadConfig != nil {
		if loadConfig.Disable {
			logx.Console().Warn("internal component load disabled", zap.String("name", name))
		}

		if loadConfig.StoreKey != "" {
			storeKey = loadConfig.StoreKey
			logx.Console().Info("internal component load config store key", zap.String("store_key", storeKey))
		}

		if loadConfig.ConfigKey != "" {
			configKey = loadConfig.ConfigKey
			logx.Console().Info("internal component load config key", zap.String("config_key", configKey))
		}
	}

	store := configx.GetStore(storeKey)
	if store == nil {
		return
	}

	l.register(name, func() error {
		return launch(storeKey, configKey)
	})
}

func (l *ComponentsLauncher) register(name string, fn func() error) {
	l.components[name] = fn
}
