package rearer

import "context"

type config struct {
	skip    int
	message string
	ctx     context.Context
}

type Option func(c *config)

func WithSkip(skip int) Option {
	return func(c *config) {
		c.skip += skip
	}
}

func WithMessage(message string) Option {
	return func(c *config) {
		c.message = message
	}
}

func WithCtx(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

func incrSkip(options []Option) []Option {
	option := []Option{
		WithSkip(1),
	}

	return append(option, options...)
}
