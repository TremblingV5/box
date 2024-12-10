package dbtx

import (
	"context"
	"sync"
)

//go:generate mockgen -source=tx.go -destination=tx_mock.go -package=dbtx TX
type TX interface {
	Commit() error
	Rollback() error
}

var (
	getTxOp    func() TX
	txInitOnce sync.Once
)

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
func doCheckTx[T any, R any](ctx context.Context, op func(tx T) (R, error)) (R, error) {
	_, getTxErr := Tx[T](ctx)
	var txCtx context.Context
	var newTx TX
	if getTxErr != nil {
		newTx = getTxOp()
		txCtx = context.WithValue(ctx, txKey{}, newTx)
	}

	r, err := op(newTx.(T))
	if getTxErr != nil {
		persist(txCtx, err)
	}

	return r, err
}

func TxDo[T any](ctx context.Context, op func(tx T) error) error {
	_, err := doCheckTx(ctx, func(tx T) (any, error) {
		return nil, op(tx)
	})
	return err
}

func TxDoGetValue[T any, R any](ctx context.Context, op func(tx T) (R, error)) (r R, err error) {
	return doCheckTx(ctx, op)
}

func TxDoGetSlice[T any, R any](ctx context.Context, op func(tx T) ([]R, error)) ([]R, error) {
	return doCheckTx(ctx, op)
}

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

func ReplaceTxPersist(ctx context.Context) (context.Context, func(err error)) {
	return withTxPersist(ctx, true, nil)
}

func ReplaceTxPersistCustom(ctx context.Context, getTx func() TX) (context.Context, func(err error)) {
	return withTxPersist(ctx, true, getTx)
}

func WithTXPersist(ctx context.Context) (context.Context, func(err error)) {
	return withTxPersist(ctx, false, nil)
}

func WithTxPersistCustom(ctx context.Context, getTx func() TX) (context.Context, func(err error)) {
	return withTxPersist(ctx, false, getTx)
}

func withTxPersist(ctx context.Context, forceReplace bool, forceGetTx func() TX) (context.Context, func(err error)) {
	checkTxErr := TxDo(ctx, func(tx TX) error {
		return nil
	})
	if checkTxErr == nil && !forceReplace {
		return ctx, func(ctx context.Context) func(err error) {
			return func(err error) {
				persist(ctx, err)
			}
		}(ctx)
	}

	var tx TX

	if forceGetTx != nil {
		tx = forceGetTx()
	} else {
		tx = getTxOp()
	}

	ctx = WithTx(ctx, tx)
	return ctx, func(ctx context.Context) func(err error) {
		return func(err error) {
			persist(ctx, err)
		}
	}(ctx)
}
