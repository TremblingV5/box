package launcher

import (
	"flag"
	"github.com/TremblingV5/box/configx"
	"os"
	"time"
)

type Launcher struct {
	flagSet *flag.FlagSet
	config  *configx.Config
	watcher []configx.WatchFunc

	stopTimeout time.Duration
}

func setFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("box-launcher", flag.ExitOnError)
	flagSet.String("config", "config/config.yaml", "config file path")
	return flagSet
}

func New() *Launcher {
	return &Launcher{
		flagSet:     setFlagSet(),
		stopTimeout: 30 * time.Second,
	}
}

func (l *Launcher) SetAppName(appName string) *Launcher {
	_ = os.Setenv("APP_NAME", appName)
	return l
}

func (l *Launcher) SetBizConfig(model any, initializer func())
