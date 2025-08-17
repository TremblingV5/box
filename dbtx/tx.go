package dbtx

import (
	"context"
	"sync"
)

//go:generate mockgen -source=tx.go -destination=tx_mock.go -package=dbtx TX

// TX defines the interface for database transaction operations.
// It provides methods to commit or rollback a transaction.
type TX interface {
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
}

var (
	getTxOp    func() TX
	txInitOnce sync.Once
)

// Init initializes the transaction getter function with a singleton pattern.
// It ensures that the transaction getter function is only set once.
func Init(op func() TX) {
	txInitOnce.Do(func() {
		getTxOp = op
	})
}

type txKey struct{}

// WithTx binds a tx to the context
func WithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Tx gets and returns the tx from the context
func Tx[T any](ctx context.Context) (tx T, err error) {
	v := ctx.Value(txKey{})
	if v == nil {
		return tx, ErrNotBindWithTx
	}

	tx, ok := v.(T)
	if !ok {
		return tx, ErrNotBindWithTheSpecifiedTx
	}

	return tx, nil
}

// doCheckTx checks the tx and executes the operation
// If there's no transaction in the context, it creates a new one
// After executing the operation, it persists the transaction (commit/rollback)
func doCheckTx[T any, R any](ctx context.Context, op func(tx T) (R, error)) (R, error) {
	tx, getTxErr := Tx[T](ctx)
	var txCtx context.Context
	var newTx TX
	if getTxErr != nil {
		newTx = getTxOp()
		txCtx = context.WithValue(ctx, txKey{}, newTx)
	}

	var r R
	var err error
	if getTxErr != nil {
		r, err = op(newTx.(T))
	} else {
		r, err = op(tx)
	}

	if getTxErr != nil {
		persist(txCtx, err)
	}

	return r, err
}

// TxDo executes a transaction operation that doesn't return a value
// It's useful for operations that only need to perform actions without returning data
func TxDo[T any](ctx context.Context, op func(tx T) error) error {
	_, err := doCheckTx(ctx, func(tx T) (any, error) {
		return nil, op(tx)
	})
	return err
}

// TxDoGetValue executes a transaction operation that returns a single value
// It's useful for operations like SELECT queries that return a single result
func TxDoGetValue[T any, R any](ctx context.Context, op func(tx T) (R, error)) (r R, err error) {
	return doCheckTx(ctx, op)
}

// TxDoGetSlice executes a transaction operation that returns a slice of values
// It's useful for operations like SELECT queries that return multiple results
func TxDoGetSlice[T any, R any](ctx context.Context, op func(tx T) ([]R, error)) ([]R, error) {
	return doCheckTx(ctx, op)
}

// persist handles the commit or rollback of a transaction based on the error status
// If err is nil, it attempts to commit the transaction
// If err is not nil or commit fails, it attempts to rollback the transaction
func persist(ctx context.Context, err error) (error, error) {
	if err == nil {
		if commitErr := TxDo(ctx, func(tx TX) error {
			return tx.Commit()
		}); commitErr == nil {
			return err, nil
		}
	}

	if rollbackErr := TxDo(ctx, func(tx TX) error {
		return tx.Rollback()
	}); rollbackErr != nil {
		return err, rollbackErr
	}

	return err, nil
}

// ReplaceTxPersist creates a new transaction context and returns a function to persist it
// It forces replacement of any existing transaction in the context
func ReplaceTxPersist(ctx context.Context) (context.Context, func(err error)) {
	return withTxPersist(ctx, true, nil)
}

// ReplaceTxPersistCustom creates a new transaction context with a custom transaction getter and returns a function to persist it
// It forces replacement of any existing transaction in the context
func ReplaceTxPersistCustom(ctx context.Context, getTx func() TX) (context.Context, func(err error)) {
	return withTxPersist(ctx, true, getTx)
}

// WithTXPersist creates a transaction context if one doesn't exist and returns a function to persist it
// If a transaction already exists in the context and forceReplace is false, it uses the existing one
func WithTXPersist(ctx context.Context) (context.Context, func(err error)) {
	return withTxPersist(ctx, false, nil)
}

// WithTxPersistCustom creates a transaction context with a custom transaction getter if one doesn't exist and returns a function to persist it
// If a transaction already exists in the context and forceReplace is false, it uses the existing one
func WithTxPersistCustom(ctx context.Context, getTx func() TX) (context.Context, func(err error)) {
	return withTxPersist(ctx, false, getTx)
}

// withTxPersist is the core implementation for transaction persistence handling
// It either uses an existing transaction or creates a new one based on the parameters
func withTxPersist(ctx context.Context, forceReplace bool, forceGetTx func() TX) (context.Context, func(err error)) {
	// Check if there's already a transaction in the context
	// We only check if a transaction is bound, without actually performing any transaction operations
	_, checkTxErr := Tx[TX](ctx)

	// If there's no error (meaning there is a transaction) and we're not forcing replacement,
	// return the same context with a persist function
	if checkTxErr == nil && !forceReplace {
		return ctx, func(err error) {
			persist(ctx, err)
		}
	}

	// Otherwise, create a new transaction
	var tx TX

	if forceGetTx != nil {
		tx = forceGetTx()
	} else {
		tx = getTxOp()
	}

	ctx = WithTx(ctx, tx)
	return ctx, func(err error) {
		persist(ctx, err)
	}
}
