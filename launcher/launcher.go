package launcher

import (
	"flag"
	"fmt"
	"github.com/TremblingV5/box/internal/shutdown"
	"github.com/TremblingV5/box/logx"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/TremblingV5/box/configx"
)

type Launcher struct {
	flagSet *flag.FlagSet
	config  *configx.Config
	watcher []configx.WatchFunc

	i18nPath string

	beforeConfigInitHandlers  []func()
	afterConfigInitHandlers   []func()
	beforeServerStartHandlers []func()
	afterServerStartHandlers  []func()
	shutdownHandlers          []func()

	stopTimeout time.Duration

	servers []Server
}

func setFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("box-launcher", flag.ExitOnError)
	flagSet.String("config", "config/config.yaml", "config file path")
	return flagSet
}

func New() *Launcher {
	atomicLevel := zap.NewAtomicLevel()
	loggerFactory := logx.NewFactory(zap.NewNop(), &atomicLevel)
	logx.SetGlobalFactory(loggerFactory)

	return &Launcher{
		flagSet:     setFlagSet(),
		stopTimeout: 30 * time.Second,
		config: &configx.Config{
			SubConfig:        make(map[string]configx.SubConfig),
			ComponentLoadMap: make(configx.ComponentLoadMap),
		},
	}
}

func (l *Launcher) SetAppName(appName string) *Launcher {
	_ = os.Setenv("APP_NAME", appName)
	return l
}

func (l *Launcher) SetBizConfig(model any, initializer func()) {

}

func (l *Launcher) WatchConfig(watchFunctions ...configx.WatchFunc) *Launcher {
	l.watcher = append(l.watcher, watchFunctions...)
	return l
}

func (l *Launcher) AddBeforeConfigInitHandler(handlers ...func()) *Launcher {
	l.beforeConfigInitHandlers = append(l.beforeConfigInitHandlers, handlers...)
	return l
}

func (l *Launcher) AddAfterConfigInitHandler(handlers ...func()) *Launcher {
	l.afterConfigInitHandlers = append(l.afterConfigInitHandlers, handlers...)
	return l
}

func (l *Launcher) AddBeforeServerStartHandler(handlers ...func()) *Launcher {
	l.beforeServerStartHandlers = append(l.beforeServerStartHandlers, handlers...)
	return l
}

func (l *Launcher) AddAfterServerStartHandler(handlers ...func()) *Launcher {
	l.afterServerStartHandlers = append(l.afterServerStartHandlers, handlers...)
	return l
}

func (l *Launcher) AddShutdownHandler(handlers ...func()) *Launcher {
	l.shutdownHandlers = append(l.shutdownHandlers, handlers...)
	return l
}

func (l *Launcher) SetI18NPath(path string) *Launcher {
	l.i18nPath = path
	return l
}

func (l *Launcher) Run(options ...Option) {
	dir, _ := os.Getwd()
	logx.Console().Info(fmt.Sprintf("current work directory: %s", dir))

	l.runBeforeConfigInitHandlers()
	configx.Load(l.config)
	logx.Console().Info("start to launch")
	l.runServer()
}

func (l *Launcher) RunWithoutServer(options ...Option) {

}

func (l *Launcher) runServer() {
	l.runAfterConfigInitHandlers()
	l.launchComponents(l.config.ComponentLoadMap)
	l.runBeforeServerStartHandlers()

	<-StartAll(l.servers...)

	l.runAfterServerStartHandlers()

	<-shutdown.FiredCh()
	shutdown.Wait(l.stopTimeout)
}

func (l *Launcher) runHandlers(handlers []func()) {
	for _, handler := range handlers {
		handler()
	}
}

func (l *Launcher) runBeforeConfigInitHandlers() {
	logx.Console().Info("start to run handlers before config init")
	l.runHandlers(l.beforeConfigInitHandlers)
	logx.Console().Info("finish to run handlers before config init")
}

func (l *Launcher) runAfterConfigInitHandlers() {
	logx.Console().Info("start to run handlers after config init")
	l.runHandlers(l.afterConfigInitHandlers)
	logx.Console().Info("finish to run handlers after config init")
}

func (l *Launcher) runBeforeServerStartHandlers() {
	logx.Console().Info("start to run handlers before server start")
	l.runHandlers(l.beforeServerStartHandlers)
	logx.Console().Info("finish to run handlers before server start")
}

func (l *Launcher) runAfterServerStartHandlers() {
	logx.Console().Info("start to run handlers after server start")
	l.runHandlers(l.afterServerStartHandlers)
	logx.Console().Info("finish to run handlers after server start")
}

func (l *Launcher) runShutdownHandlers() {
	logx.Console().Info("start to run shutdown handlers")
	l.runHandlers(l.shutdownHandlers)
	logx.Console().Info("finish to run shutdown handlers")
}

func (l *Launcher) launchComponents(loadMap configx.ComponentLoadMap) {
	launcher := newComponentsLauncher(loadMap)
	launcher.launch()
}

func (l *Launcher) AddServer(server ...Server) {
	l.servers = append(l.servers, server...)
}
