package main

import (
	"context"
	"log"

	"github.com/TremblingV5/box/dbtx"
)

type fakeTx struct {
	name string
}

func (t *fakeTx) Commit() error {
	log.Println("commit")
	return nil
}

func (t *fakeTx) Rollback() error {
	log.Println("rollback")
	return nil
}

func init() {
	dbtx.Init(func() dbtx.TX {
		return &fakeTx{
			name: "fakeTx 001",
		}
	})
}

func serviceMethod(ctx context.Context) (err error) {
	ctx, persist := dbtx.WithTXPersist(ctx)
	defer func() {
		persist(err)
	}()

	repositoryMethod(ctx)

	return nil
}

func repositoryMethod(ctx context.Context) error {
	return dbtx.TxDo(ctx, func(tx *fakeTx) error {
		log.Println("get tx:", tx)
		return nil
	})
}

func main() {
	serviceMethod(context.Background())
}
