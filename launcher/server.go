package launcher

import (
	"github.com/TremblingV5/box/gofer"
)

type Server interface {
	Start() <-chan struct{}
}

var globalClosedChan = make(chan struct{})

func init() {
	close(globalClosedChan)
}

func StartAll(servers ...Server) chan struct{} {
	if len(servers) == 0 {
		return globalClosedChan
	}

	var chs []<-chan struct{}

	for _, server := range servers {
		chs = append(chs, server.Start())
	}

	serverReadyCh := make(chan struct{})

	gofer.Go(func() {
		for _, ch := range chs {
			<-ch
		}

		close(serverReadyCh)
	})

	return serverReadyCh
}
