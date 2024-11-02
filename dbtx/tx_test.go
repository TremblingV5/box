package dbtx

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

// Test_WithTx tests WithTx function
func Test_WithTx(t *testing.T) {
	t.Run("WithTx should bind a tx to the context", func(t *testing.T) {
		tx := &MockTX{}
		ctx := context.Background()
		ctx = WithTx(ctx, tx)
		require.NotNil(t, ctx.Value(txKey{}))
		require.Equal(t, tx, ctx.Value(txKey{}))
	})
}

// Test_Tx tests Tx function
func Test_Tx(t *testing.T) {
	t.Run("Tx should get and return the tx from the context", func(t *testing.T) {
		tx := &MockTX{}
		ctx := context.WithValue(context.Background(), txKey{}, tx)
		txFromCtx, err := Tx[*MockTX](ctx)
		require.Nil(t, err)
		require.Equal(t, tx, txFromCtx)
	})

	t.Run("Tx should return an error if tx not bind to ctx", func(t *testing.T) {
		ctx := context.Background()
		txFromCtx, err := Tx[*MockTX](ctx)
		require.NotNil(t, err)
		require.Equal(t, ErrNotBindWithTx, err)
		require.Nil(t, txFromCtx)
	})

	t.Run("Tx should return an error if tx not bind with the correct ctx", func(t *testing.T) {
		tx := struct{}{}
		ctx := context.WithValue(context.Background(), txKey{}, tx)
		txFromCtx, err := Tx[*MockTX](ctx)
		require.NotNil(t, err)
		require.Equal(t, ErrNotBindWithTheSpecifiedTx, err)
		require.Nil(t, txFromCtx)
	})
}

func Test_doCheckTx(t *testing.T) {
	tx := NewMockTX(gomock.NewController(t))

	t.Run("normal", func(t *testing.T) {
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		ctx = WithTx(ctx, tx)
		r, err := doCheckTx(ctx, func(tx *MockTX) (string, error) {
			return "result", nil
		})
		require.NotNil(t, r)
		require.Equal(t, "result", r)
		require.Nil(t, err)
	})

	t.Run("tx not bind to the context", func(t *testing.T) {
		tx.EXPECT().Commit().Return(nil)
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		r, err := doCheckTx(ctx, func(tx *MockTX) (string, error) {
			return "result", nil
		})
		require.NotNil(t, r)
		require.Equal(t, "result", r)
		require.Nil(t, err)
	})

	t.Run("op return an error", func(t *testing.T) {
		tx.EXPECT().Rollback().Return(nil)
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		_, err := doCheckTx(ctx, func(tx *MockTX) (string, error) {
			return "", errors.New("some error")
		})
		require.NotNil(t, err)
		require.Equal(t, "some error", err.Error())
	})
}

func Test_persist(t *testing.T) {
	tx := NewMockTX(gomock.NewController(t))

	t.Run("no error", func(t *testing.T) {
		tx.EXPECT().Commit().Return(nil)
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		ctx = WithTx(ctx, tx)
		err, commitErr := persist(ctx, nil)
		require.Nil(t, err)
		require.Nil(t, commitErr)
	})

	t.Run("failed to commit", func(t *testing.T) {
		tx.EXPECT().Commit().Return(errors.New("commit error"))
		tx.EXPECT().Rollback().Return(nil)
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		ctx = WithTx(ctx, tx)
		err, rollbackErr := persist(ctx, nil)
		require.Nil(t, err)
		require.Nil(t, rollbackErr)
	})

	t.Run("failed to commit, failed to rollback", func(t *testing.T) {
		tx.EXPECT().Commit().Return(errors.New("commit error"))
		tx.EXPECT().Rollback().Return(errors.New("rollback error"))
		Init(func() TX {
			return tx
		})

		ctx := context.Background()
		ctx = WithTx(ctx, tx)
		err, rollbackErr := persist(ctx, nil)
		require.Nil(t, err)
		require.NotNil(t, rollbackErr)
		require.Equal(t, "rollback error", rollbackErr.Error())
	})
}
