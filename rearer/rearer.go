package rearer

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func Recover(options ...Option) {
	err := recover()
	if err != nil {
		return
	}

	LogRecoverStack(err, options...)
}

func RecoverWithCtx(ctx context.Context, options ...Option) {
	options = append(options, WithCtx(ctx))
	Recover(options...)
}

func LogRecoverStack(err any, options ...Option) {
	logStack(err, incrSkip(options)...)
}

func logStack(err any, options ...Option) {
	var caller string

	cfg := &config{
		skip: 3,
		ctx:  context.Background(),
	}

	for _, option := range options {
		option(cfg)
	}

	if pc, file, line, ok := runtime.Caller(cfg.skip); ok {
		caller = fmt.Sprintf("func: %s, file: %s:%d", runtime.FuncForPC(pc), file, line)
	}

	stack := string(debug.Stack())

	hookInfo := &HookInfo{
		Context:    cfg.ctx,
		PanicError: err,
		Caller:     caller,
		Stack:      stack,
		Message:    cfg.message,
	}

	for _, hook := range globalHooks {
		hook(hookInfo)
	}

	var messageBuilder strings.Builder

	messageBuilder.WriteString("[Recovery from panic]")
	if cfg.message != "" {
		messageBuilder.WriteString(fmt.Sprintf(" message: %s", cfg.message))
	}

	// TODO: 写日志
}
