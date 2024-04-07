package configx

import (
	"sync"
	"sync/atomic"
)

var (
	globalStoreMap  sync.Map
	globalStoreInit atomic.Bool
)

const (
	StoreKeyDefault = "default"
)

func GetStore(storeKey string) Store {
	if store, ok := globalStoreMap.Load(storeKey); ok {
		return store.(Store)
	}

	return nil
}
