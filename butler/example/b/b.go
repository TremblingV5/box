package b

import (
	"log"

	"github.com/TremblingV5/box/butler"
	"github.com/TremblingV5/box/butler/example/c"
)

func init() {
	butler.Provide(func() IB {
		return &B{
			c:    butler.Invoke[c.IC](),
			Name: "real B",
		}
	})
}

type IB interface {
	DoB()
}

type B struct {
	c    c.IC
	Name string
}

func (b *B) DoB() {
	log.Println("I'm B, I have a C. And C's name is", b.c.(*c.C).Name)

	b.c.DoC()
}
