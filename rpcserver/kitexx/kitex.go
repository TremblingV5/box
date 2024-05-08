package kitexx

import (
	"github.com/TremblingV5/box/gofer"
	"github.com/cloudwego/kitex/server"
	"sync"
)

type KitexServer[T any] struct {
	initMethod func(T, ...server.Option) server.Server
	handler    T
	options    []server.Option
}

func NewKitexServer[T any](
	initMethod func(T, ...server.Option) server.Server,
	handler T,
	options ...server.Option,
) *KitexServer[T] {
	return &KitexServer[T]{
		initMethod: initMethod,
		handler:    handler,
		options:    options,
	}
}

func (s *KitexServer[T]) Start() <-chan struct{} {
	readyChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(1)

	gofer.Go(func() {
		wg.Done()

		server := s.initMethod(s.handler, s.options...)
		if err := server.Run(); err != nil {
			panic(err)
		}
	})

	gofer.Go(func() {
		wg.Wait()

		close(readyChan)
	})

	return readyChan
}
