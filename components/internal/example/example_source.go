package example

//
//import (
//	"github.com/TremblingV5/box/components"
//	"github.com/TremblingV5/box/configx"
//	"github.com/TremblingV5/box/launcher"
//)
//
//// InitComponentConfig only used to init config for example code, do not use it in other project
//func InitComponentConfig[T any](
//	rootConfigPath, componentName, componentConfigKey string,
//	initMethod func(cfg components.ConfigMap[T]) error,
//) {
//	configx.SetRootConfigPath(rootConfigPath)
//	configx.ResolveWatcher(
//		configx.Watch("database", configx.WithModel(&configx.Config{})),
//	)
//	configx.Init()
//
//	store := configx.GetStore(configx.StoreKeyConfig)
//	loadMap := make(configx.ComponentLoadMap)
//	store.UnmarshalKey("componentLoad", &loadMap)
//
//	l := launcher.NewComponentsLauncher(loadMap)
//	l.CustomLaunch(func(l *launcher.ComponentsLauncher) {
//		l.LaunchComponent(componentName, componentConfigKey, func(storeKey, configKey string) error {
//			return components.Load(storeKey, configKey, initMethod).Start()
//		})
//	})
//}
