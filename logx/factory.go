package logx

import "sync"

var (
	globalMutex   sync.RWMutex
	globalFactory *Factory
)

type Factory struct {
	logger     *logger
	logOptions *logOptions
}

func GetGlobalFactory() *Factory {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	return globalFactory
}

func SetGlobalFactory(factory *Factory) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	globalFactory = factory
}
