package main

import (
	"fmt"
	"github.com/TremblingV5/box/configx"
)

type BizConfig struct {
	Test struct {
		Test1 int `yaml:"test1" json:"test1" mapstructure:"test1"`
		Test2 int `yaml:"test2" json:"test2" mapstructure:"test2"`
	} `yaml:"test" json:"test" mapstructure:"test"`
}

const (
	storeKeyBiz  = "biz"
	storeKeyBiz1 = "biz1"
)

func init() {
	configx.SetRootConfigPath("configx/example/file/config")
}

// GetBizConfig should be defined in config/config.go in real project
func GetBizConfig() *BizConfig {
	return configx.LoadCustomModel[*BizConfig](storeKeyBiz)
}

// GetBiz1Config should be defined in config/config.go in real project
func GetBiz1Config() *BizConfig {
	return configx.LoadCustomModel[*BizConfig](storeKeyBiz1)
}

func main() {
	configx.SetWatchMap(
		configx.Watch(storeKeyBiz, configx.WithModel(&BizConfig{})),
		configx.Watch(storeKeyBiz1, configx.WithModel(&BizConfig{})),
	)
	configx.Init()

	bizConfig := GetBizConfig()
	fmt.Println(bizConfig.Test.Test1)
	fmt.Println(bizConfig.Test.Test2)

	biz1Config := GetBiz1Config()
	fmt.Println(biz1Config.Test.Test1)
	fmt.Println(biz1Config.Test.Test2)
}
