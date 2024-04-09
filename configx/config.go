package configx

var (
	ConfigTypeEnvKey = "CONFIG_TYPE"
	RootConfigPath   = "config/local"
)

func SetConfigTypeEnvKey(key string) {
	ConfigTypeEnvKey = key
}

func SetRootConfigPath(path string) {
	RootConfigPath = path
}

type Config struct {
	System           *SystemConfig        `yaml:"system" json:"system" mapstructure:"system"`
	SubConfig        map[string]SubConfig `yaml:"subConfig" json:"subConfig" mapstructure:"subConfig"`
	ComponentLoadMap ComponentLoadMap     `yaml:"componentLoadMap" json:"componentLoadMap" mapstructure:"componentLoadMap"`
}

type SystemConfig struct {
	AppName      string `yaml:"appName" json:"appName" mapstructure:"appName"`
	RouterPrefix string `yaml:"routerPrefix" json:"routerPrefix" mapstructure:"routerPrefix"`
	Host         string `yaml:"host" json:"host" mapstructure:"host"`
	HttpPort     int    `yaml:"httpPort" json:"httpPort" mapstructure:"httpPort"`
	GrpcPort     int    `yaml:"grpcPort" json:"grpcPort" mapstructure:"grpcPort"`
	Mode         string `yaml:"mode" json:"mode" mapstructure:"mode"`
	Env          string `yaml:"env" json:"env" mapstructure:"env"`
	EnablePprof  bool   `yaml:"enablePprof"	json:"enablePprof" mapstructure:"enablePprof"`
}

type SubConfig struct {
	Type     string `yaml:"type" json:"type" mapstructure:"type"` // inherit | file | viper
	Location string `yaml:"location" json:"location" mapstructure:"location"`
}

type TraceConfig struct {
	Enable  bool   `yaml:"enable" json:"enable" mapstructure:"enable"`
	Type    string `yaml:"type" json:"type" mapstructure:"type"`
	Address string `yaml:"address" json:"address" mapstructure:"address"`
}

type ComponentLoadMap map[string]*ComponentLoadConfig

type ComponentLoadConfig struct {
	Disable   bool   `yaml:"disable" json:"disable" mapstructure:"disable"`
	StoreKey  string `yaml:"storeKey" json:"storeKey" mapstructure:"storeKey"`
	ConfigKey string `yaml:"configKey" json:"configKey" mapstructure:"configKey"`
}
