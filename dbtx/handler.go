package dbtx

type Handler struct {
	getTxOp func() TX
}

func New(op func() TX) *Handler {
	return &Handler{
		getTxOp: op,
	}
}

func (h *Handler) WithTXPersist(ctx context.Context) (context.Context, func(err error)) {
	return WithTxPersistCustom(ctx, h.getTxOp)
}

func (h *Handler) ReplaceTxPersist(ctx context.Context) (context.Context, func(err error)) {
	return ReplaceTxPersistCustom(ctx, h.getTxOp)
}
