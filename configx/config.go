package configx

type Config struct {
	System           *SystemConfig        `yaml:"system"`
	SubConfig        map[string]SubConfig `yaml:"subConfig"`
	ComponentLoadMap ComponentLoadMap     `yaml:"componentLoadMap"`
}

type SystemConfig struct {
	AppName      string `yaml:"appName"`
	RouterPrefix string `yaml:"routerPrefix"`
	Host         string `yaml:"host"`
	HttpPort     int    `yaml:"httpPort"`
	GrpcPort     int    `yaml:"grpcPort"`
	Mode         string `yaml:"mode"`
	Env          string `yaml:"env"`
	EnablePprof  bool   `yaml:"enablePprof"`
}

type SubConfig struct {
	Type     string `yaml:"type"` // inherit | file | viper
	Location string `yaml:"location"`
}

type TraceConfig struct {
	Enable  bool   `yaml:"enable"`
	Type    string `yaml:"type"`
	Address string `yaml:"address"`
}

type ComponentLoadMap map[string]*ComponentLoadConfig

type ComponentLoadConfig struct {
	Disable   bool   `yaml:"disable"`
	StoreKey  string `yaml:"storeKey"`
	ConfigKey string `yaml:"configKey"`
}
