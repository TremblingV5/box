package main

import (
	"context"
	"log"

	"github.com/TremblingV5/box/dbtx"
)

// fakeTx is a mock implementation of a database transaction
type fakeTx struct {
	name string
}

// Commit simulates committing a transaction
func (t *fakeTx) Commit() error {
	log.Println("commit")
	return nil
}

// Rollback simulates rolling back a transaction
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

// serviceMethod demonstrates a typical service method that uses database transactions
// It shows how to properly handle transaction persistence with defer
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

// repositoryMethod1 demonstrates how to retrieve a transaction from context and use it
func repositoryMethod1(ctx context.Context) error {
	// you can ignore this error only you can ensure that the tx is always available
	tx, _ := dbtx.Tx[*fakeTx](ctx)
	log.Println("get tx and do something:", tx)
	return nil
}

// repositoryMethod2 demonstrates using TxDo for operations that don't return values
func repositoryMethod2(ctx context.Context) error {
	// dbtx.TxDo used to do something without returning value.
	return dbtx.TxDo(ctx, func(tx *fakeTx) error {
		log.Println("get tx and do something in dbtx.TxDo:", tx)
		return nil
	})
}

// repositoryMethod3 demonstrates using TxDoGetValue for operations that return a single value
func repositoryMethod3(ctx context.Context) (string, error) {
	// dbtx.TxDoGetValue used to do something and return a single value.
	return dbtx.TxDoGetValue(ctx, func(tx *fakeTx) (string, error) {
		log.Println("get tx and do something in dbtx.TxDoGetValue:", tx)
		return "tx do get value successfully", nil
	})
}

// repositoryMethod4 demonstrates using TxDoGetSlice for operations that return a slice
func repositoryMethod4(ctx context.Context) ([]string, error) {
	// dbtx.TxDoGetSlice used to do something and return a slice.
	return dbtx.TxDoGetSlice(ctx, func(tx *fakeTx) ([]string, error) {
		log.Println("get tx and do something in dbtx.TxDoGetSlice:", tx)
		return []string{"tx", "do", "get", "slice", "successfully"}, nil
	})

}

// main executes the service method to demonstrate transaction handling
func main() {
	serviceMethod(context.Background())
}
