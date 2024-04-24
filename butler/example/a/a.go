package a

import (
	"log"

	"github.com/TremblingV5/box/butler"
	"github.com/TremblingV5/box/butler/example/b"
)

func init() {
	butler.Provide(func() IA {
		return &A{
			b: butler.Invoke[b.IB](),
		}
	})
}

type IA interface {
	DoA()
}

type A struct {
	b b.IB
}

func (a *A) DoA() {
	log.Println("I'm A, I have a B. And B's name is", a.b.(*b.B).Name)

	a.b.DoB()
}
