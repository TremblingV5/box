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
	// Used to define how we get the `tx`. It's useful when we have
	// more than one `tx` in the project.
	dbtx.Init(func() dbtx.TX {
		return &fakeTx{
			name: "fakeTx 001",
		}
	})
}

func serviceMethod(ctx context.Context) (err error) {
	// use the following four lines to add a `tx` to context and ensure transaction can be committed.
	ctx, persist := dbtx.WithTXPersist(ctx)
	defer func() {
		persist(err)
	}()

	if err := repositoryMethod1(ctx); err != nil {
		return err
	}

	if err := repositoryMethod2(ctx); err != nil {
		return err
	}

	result1, err := repositoryMethod3(ctx)
	if err != nil {
		return err
	}

	log.Println("repositoryMethod3 result:", result1)

	result2, err := repositoryMethod4(ctx)
	if err != nil {
		return err
	}

	for index, item := range result2 {
		log.Println("repositoryMethod4 result:", index, item)
	}

	return nil
}

func repositoryMethod1(ctx context.Context) error {
	// you can ignore this error only you can ensure that the tx is always available
	tx, _ := dbtx.Tx[*fakeTx](ctx)
	log.Println("get tx and do something:", tx)
	return nil
}

func repositoryMethod2(ctx context.Context) error {
	// dbtx.TxDo used to do something without returning value.
	return dbtx.TxDo(ctx, func(tx *fakeTx) error {
		log.Println("get tx and do something in dbtx.TxDo:", tx)
		return nil
	})
}

func repositoryMethod3(ctx context.Context) (string, error) {
	// dbtx.TxDoGetValue used to do something and return a single value.
	return dbtx.TxDoGetValue(ctx, func(tx *fakeTx) (string, error) {
		log.Println("get tx and do something in dbtx.TxDoGetValue:", tx)
		return "tx do get value successfully", nil
	})
}

func repositoryMethod4(ctx context.Context) ([]string, error) {
	// dbtx.TxDoGetSlice used to do something and return a slice.
	return dbtx.TxDoGetSlice(ctx, func(tx *fakeTx) ([]string, error) {
		log.Println("get tx and do something in dbtx.TxDoGetSlice:", tx)
		return []string{"tx", "do", "get", "slice", "successfully"}, nil
	})

}

func main() {
	serviceMethod(context.Background())
}
