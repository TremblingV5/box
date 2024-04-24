package c

import (
	"log"

	"github.com/TremblingV5/box/butler"
)

func init() {
	butler.Provide(func() IC {
		return &C{
			Name: "real C",
		}
	})
}

type IC interface {
	DoC()
}

type C struct {
	Name string
}

func (c *C) DoC() {
	log.Println("I'm C")
}
