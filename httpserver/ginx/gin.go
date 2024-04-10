package ginx

import (
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/TremblingV5/box/gofer"
	"github.com/TremblingV5/box/logx"
	"github.com/gin-gonic/gin"
)

type RegisterGinRouter func(*gin.RouterGroup)

type GinServer struct {
	addr        string
	listener    net.Listener
	server      *http.Server
	engine      *gin.Engine
	group       *gin.RouterGroup
	ginMode     string
	contextPath string
}

func newGinEngine() *gin.Engine {
	engine := gin.New()
	engine.ContextWithFallback = true
	return engine
}

func NewGinServer(addr, ginMode, contextPath string) *GinServer {
	ginServer := &GinServer{
		addr:        addr,
		server:      defaultHttpServer(),
		engine:      newGinEngine(),
		ginMode:     ginMode,
		contextPath: contextPath,
	}

	if ginServer.ginMode == gin.DebugMode {
		ginServer.engine.Use(gin.Logger())
	}

	ginServer.group = ginServer.engine.Group(ginServer.contextPath)
	ginServer.server.Handler = ginServer.engine

	return ginServer
}

func (s *GinServer) Start() <-chan struct{} {
	if s.listener == nil {
		listener, err := net.Listen("tcp", s.addr)
		if err != nil {
			logx.Console().Fatal("failed to listen http addr")
		}
		s.listener = listener
	}

	readyChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	gofer.Go(func() {
		wg.Done()

		err := s.server.Serve(s.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logx.Fatal("failed to start http server")
		}
	})

	gofer.Go(func() {
		wg.Wait()

		listenAddr := s.listener.Addr().String()
		logx.Info("starting http server: %s", listenAddr)

		close(readyChan)
	})

	return readyChan
}

func (s *GinServer) RegisterHttpHandlers(fns ...RegisterGinRouter) {
	for _, fn := range fns {
		fn(s.group)
	}
}
