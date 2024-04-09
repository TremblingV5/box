package configx

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	_ "github.com/spf13/viper/remote"
)

var (
	globalStoreMap  sync.Map
	globalStoreInit atomic.Bool
	globalWatchMap  = make(WatchMap)
)

const (
	ConfigCenterUrlEnvKey      = "CONFIG_CENTER_URL"
	ConfigCenterStoreKeyPrefix = "CONFIG_CENTER_STORE_KEY_PREFIX/"
	StoreKeyDefault            = "default"
	StoreKeyConfig             = "config"
)

func GetStore(storeKey string) *Store {
	if store, ok := globalStoreMap.Load(storeKey); ok {
		return store.(*Store)
	}

	return nil
}

func SetStore(storeKey string, store *Store) {
	globalStoreMap.Store(storeKey, store)
	globalStoreInit.Store(true)
}

func LoadCustomModel[T any](storeKey string) T {
	store := GetStore(storeKey)
	if store == nil {
		return *new(T)
	}

	return store.aval.Load().(T)
}

func IsInit() bool {
	return globalStoreInit.Load()
}

func SetWatchMap(watchFunctions ...WatchFunc) {
	for _, watchFunction := range watchFunctions {
		watchFunction(globalWatchMap)
	}
}

func Init() {
	if globalStoreInit.Load() {
		return
	}

	configType := os.Getenv(ConfigTypeEnvKey)
	switch configType {
	case "ETCD":
		initFromEtcd()
	default:
		initFromYaml()
	}

	solveBaseConfig()
}

func loadConfigFromFile(path, fileType string) *Store {
	store := NewStore()
	store.viper.SetConfigFile(path)
	store.viper.SetConfigType(fileType)
	store.viper.ReadInConfig()
	return store
}

func loadConfigFromEtcd(url, storeKey string) *Store {
	store := NewStore()
	store.viper.AddRemoteProvider("etcd3", url, ConfigCenterStoreKeyPrefix+storeKey)
	store.viper.SetConfigType("json")
	store.viper.ReadRemoteConfig()
	return store
}

func solveBaseConfig() {
	configCenterUrl := os.Getenv(ConfigCenterUrlEnvKey)

	baseConfigStore := GetStore(StoreKeyConfig)
	subConfigs := make(map[string]SubConfig)
	_ = baseConfigStore.UnmarshalKey("subConfig", &subConfigs)

	componentsLoadConfigs := make(ComponentLoadMap)
	_ = baseConfigStore.UnmarshalKey("componentLoad", &componentsLoadConfigs)

	for key, watchSetting := range globalWatchMap {
		subConfig, ok := subConfigs[key]
		if !ok {
			continue
		}

		store := NewStore()
		switch subConfig.Type {
		case "inherit":
			store.viper = baseConfigStore.viper
		case "file":
			subConfigPath := filepath.Join(RootConfigPath, subConfig.Location) + ".yaml"
			store.viper.SetConfigFile(subConfigPath)
			store.viper.SetConfigType("yaml")
			store.viper.ReadInConfig()
		case "viper":
			store.viper.AddRemoteProvider("etcd3", configCenterUrl, ConfigCenterStoreKeyPrefix+key)
			store.viper.SetConfigType("yaml")
			store.viper.ReadRemoteConfig()
		default:
			panic("unsupported subConfig type")
		}

		if watchSetting.model != nil {
			store.viper.Unmarshal(watchSetting.model)
			store.aval.Store(watchSetting.model)
		}

		SetStore(key, store)
	}
}

func initFromYaml() {
	totalConfigFilePath := filepath.Join(RootConfigPath, "config.yaml")
	SetStore(StoreKeyConfig, loadConfigFromFile(totalConfigFilePath, "yaml"))
}

func initFromEtcd() {
	configCenterUrl := os.Getenv(ConfigCenterUrlEnvKey)
	SetStore(StoreKeyConfig, loadConfigFromEtcd(configCenterUrl, StoreKeyConfig))
}
