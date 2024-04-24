package main

import (
	"github.com/TremblingV5/box/butler"
	"github.com/TremblingV5/box/butler/example/a"
)

func main() {
	value := butler.Invoke[a.IA]()
	value.DoA()
}
