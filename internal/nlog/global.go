package nlog

import (
	"sync"
	"sync/atomic"
)

var (
	globalLoggerMap     = sync.Map{}
	globalLoggerBuilder = atomic.Value{}
)

func SetLogger(name string, logger Logger) {
	globalLoggerMap.Store(name, logger)
}

func GetLogger(name string) Logger {
	if logger, ok := globalLoggerMap.Load(name); ok {
		return logger.(Logger)
	}

	return nil
}

func GetOrNewLogger(name string) Logger {
	if logger := GetLogger(name); logger != nil {
		return logger
	}

	if builder, ok := globalLoggerBuilder.Load().(LoggerBuilder); ok {
		logger := builder(name)
		SetLogger(name, logger)
		return logger
	}

	return nil
}

type LoggerBuilder func(name string) Logger

func SetLoggerBuilder(builder LoggerBuilder) {
	globalLoggerBuilder.Store(builder)

	globalLoggerMap.Range(func(key, _ interface{}) bool {
		name := key.(string)
		SetLogger(name, builder(name))
		return true
	})
}
