package mysqlx

type Component struct {
	err error
	cfg ConfigMap
}

func Load(storeKey, configKey string) *Component {
	return &Component{}
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

}
