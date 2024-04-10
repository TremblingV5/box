package hertzx

import (
	"net"
	"sync"

	"github.com/TremblingV5/box/gofer"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

type RegisterHertzRouter func(*route.RouterGroup)

type HertzServer struct {
	addr        string
	listener    net.Listener
	hertz       *server.Hertz
	group       *route.RouterGroup
	ginMode     string
	contextPath string
}

func newHertz() *server.Hertz {
	h := server.New()
	return h
}

func NewHertzServer(addr, contextPath string) *HertzServer {
	hertzServer := &HertzServer{
		addr:        addr,
		hertz:       newHertz(),
		contextPath: contextPath,
	}

	hertzServer.group = hertzServer.hertz.Engine.Group(hertzServer.contextPath)

	return hertzServer
}

func (s *HertzServer) Start() <-chan struct{} {
	readyChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	gofer.Go(func() {
		wg.Done()

		s.hertz.Spin()
	})

	gofer.Go(func() {
		wg.Wait()
		close(readyChan)
	})

	return readyChan
}

func (s *HertzServer) RegisterHttpHandlers(fns ...RegisterHertzRouter) {
	for _, fn := range fns {
		fn(s.group)
	}
}
