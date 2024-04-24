package butler

import (
	"sync/atomic"

	"go.uber.org/dig"
)

var globalContainer = dig.New()

var lock atomic.Bool

func FinishWork() {
	lock.Store(true)
}

func Provide(f interface{}) {
	if locked := lock.Load(); locked {
		return
	}

	if err := globalContainer.Provide(f); err != nil {
		panic(err)
	}
}

func Invoke[T any]() (t T) {
	if err := globalContainer.Invoke(func(item T) {
		t = item
	}); err != nil {
		panic(err)
	}

	return t
}
