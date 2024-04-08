package mysqlx

import "github.com/TremblingV5/box/configx"

type Component struct {
	err error
	cfg ConfigMap
}

func Load(storeKey, configKey string) *Component {
	c := &Component{}
	c.cfg, c.err = LoadConfig(storeKey, configKey)
	return c
}

func (s *Component) Start() error {
	if s.err != nil {
		return s.err
	}

	return Init(s.cfg)
}

func (s *Component) GetConfig() ConfigMap {
	return s.cfg
}

func LoadConfig(storeKey, configKey string) (ConfigMap, error) {
	var cfg ConfigMap

	err := configx.GetStore(storeKey).UnmarshalKey(configKey, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
