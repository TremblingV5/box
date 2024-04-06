package logx

import (
	"context"

	"go.uber.org/zap"
)

const (
	defaultNamespace = "content"
)

type Option func(*logOptions)

type logOptions struct {
	fieldNamespace string
	addHandler     func() []zap.Field
	ctxHandler     func(ctx context.Context) []zap.Field
	errHandler     func(err error) []zap.Field
}

func defaultOptions() *logOptions {
	return &logOptions{
		fieldNamespace: defaultNamespace,
		errHandler: func(err error) []zap.Field {
			return []zap.Field{
				zap.Error(err),
			}
		},
	}
}

func WithAddHandler(f func() []zap.Field) Option {
	return func(o *logOptions) {
		o.addHandler = f
	}
}

func WithCtxHandler(f func(ctx context.Context) []zap.Field) Option {
	return func(o *logOptions) {
		o.ctxHandler = f
	}
}

func WithErrHandler(f func(err error) []zap.Field) Option {
	return func(o *logOptions) {
		o.errHandler = f
	}
}
