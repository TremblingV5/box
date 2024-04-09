package configx

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"sync"
	"sync/atomic"
)

type PreProcessorErrorFunc func(data []byte) []byte

type StoreConfig struct {
	ctx           context.Context
	notify        chan struct{}
	logOnReceive  bool
	preProcessors []PreProcessorErrorFunc
}

func defaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		ctx:          context.Background(),
		notify:       make(chan struct{}),
		logOnReceive: true,
	}
}

type Store struct {
	cfg       *StoreConfig
	aval      atomic.Value
	started   atomic.Value
	binderMap sync.Map
	viper     viper.Viper
	configs   sync.Map
}

func NewStore() *Store {
	return &Store{
		cfg:     defaultStoreConfig(),
		aval:    atomic.Value{},
		started: atomic.Value{},
		viper:   *viper.New(),
		configs: sync.Map{},
	}
}

func (s *Store) GetSubNodes(key string) *viper.Viper {
	return s.viper.Sub(key)
}

func (s *Store) UnmarshalKey(configKey string, value any) error {
	v := s.GetSubNodes(configKey)
	if v == nil {
		return errors.New("config key not found")
	}

	return v.Unmarshal(value)
}
