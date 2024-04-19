package dbtx

import (
	"context"
	"errors"
	"sync"
)

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

func WithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func Tx[T any](ctx context.Context) (tx T, err error) {
	v := ctx.Value(txKey{})
	if v == nil {
		return tx, errors.New("context not bind with tx")
	}

	tx, ok := v.(T)
	if !ok {
		return tx, errors.New("context not bind with tx")
	}

	return tx, nil
}

func TxDo[T any](ctx context.Context, op func(tx T) error) error {
	tx, err := Tx[T](ctx)
	if err != nil {
		return err
	}

	return op(tx)
}

func TxDoGetValue[T any, R any](ctx context.Context, op func(tx T) (R, error)) (r R, err error) {
	tx, err := Tx[T](ctx)
	if err != nil {
		return r, err
	}

	return op(tx)
}

func TxDoGetSlice[T any, R any](ctx context.Context, op func(tx T) ([]R, error)) ([]R, error) {
	tx, err := Tx[T](ctx)
	if err != nil {
		return nil, err
	}

	return op(tx)
}

func persist(ctx context.Context, err error) {
	if err == nil {
		if commitErr := TxDo(ctx, func(tx TX) error {
			return tx.Commit()
		}); commitErr != nil {
			return
		}

		return
	}

	if rollbackErr := TxDo(ctx, func(tx TX) error {
		return tx.Rollback()
	}); rollbackErr != nil {
		return
	}

	return
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
		return ctx, func(err error) {

		}
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
