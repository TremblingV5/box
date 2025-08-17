package dbtx

import "context"

// Handler provides a wrapper around transaction operations
// It holds a reference to the transaction getter function
type Handler struct {
	getTxOp func() TX
}

// New creates a new Handler instance with the provided transaction getter function
func New(op func() TX) *Handler {
	return &Handler{
		getTxOp: op,
	}
}

// WithTXPersist creates a transaction context using the handler's transaction getter
// It returns the new context and a function to persist (commit/rollback) the transaction
func (h *Handler) WithTXPersist(ctx context.Context) (context.Context, func(err error)) {
	return WithTxPersistCustom(ctx, h.getTxOp)
}

// ReplaceTxPersist creates a new transaction context, replacing any existing transaction
// It returns the new context and a function to persist (commit/rollback) the transaction
func (h *Handler) ReplaceTxPersist(ctx context.Context) (context.Context, func(err error)) {
	return ReplaceTxPersistCustom(ctx, h.getTxOp)
}
