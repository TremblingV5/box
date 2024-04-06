package dbtx

import (
	"context"
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

func WithTx(ctx context.Context, tx TX) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func Tx(ctx context.Context) TX {
	v := ctx.Value(txKey{})
	if v == nil {
		return nil
	}

	tx, ok := v.(TX)
	if !ok {
		return nil
	}

	return tx
}

func persist(ctx context.Context, err error) {
	tx := Tx(ctx)

	if err == nil {
		if commitErr := tx.Commit(); commitErr != nil {
			return
		}

		return
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {

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
	t0 := Tx(ctx)
	if t0 != nil && !forceReplace {
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
