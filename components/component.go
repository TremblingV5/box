package components

import "github.com/TremblingV5/box/configx"

type ConfigMap[T any] map[string]T

type Component[T any] struct {
	err        error
	cfg        ConfigMap[T]
	initMethod func(cfg ConfigMap[T]) error
}

func Load[T any](storeKey, configKey string, initMethod func(cfg ConfigMap[T]) error) *Component[T] {
	c := &Component[T]{
		initMethod: initMethod,
	}

	err := configx.GetStore(storeKey).UnmarshalKey(configKey, &c.cfg)
	c.err = err

	return c
}

func (s *Component[T]) Start() error {
	if s.err != nil {
		return s.err
	}

	return s.initMethod(s.cfg)
}

func (s *Component[T]) GetConfig() ConfigMap[T] {
	return s.cfg
}
