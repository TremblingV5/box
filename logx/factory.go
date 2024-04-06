package logx

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

var (
	globalMutex   sync.RWMutex
	globalFactory *Factory
)

type Factory struct {
	logger     *logger
	logOptions *logOptions
}

func GetGlobalFactory() *Factory {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	return globalFactory
}

func SetGlobalFactory(factory *Factory) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	globalFactory = factory
}

func NewFactory(logger *zap.Logger, atomicLevel *zap.AtomicLevel, options ...Option) *Factory {
	if atomicLevel == nil {
		atomicLevel = GetGlobalFactory().logger.atomicLevel
	}

	return newFactory(newLogger(logger, atomicLevel), options...)
}

func newFactory(logger *logger, options ...Option) *Factory {
	defaults := defaultOptions()

	for _, option := range options {
		option(defaults)
	}

	return &Factory{
		logger:     logger,
		logOptions: defaults,
	}
}

func (f *Factory) SetOptions(options ...Option) *Factory {
	opts := *f.logOptions

	for _, option := range options {
		option(&opts)
	}

	return &Factory{
		logger:     f.logger,
		logOptions: &opts,
	}
}

func (f *Factory) GetZapLogger() *zap.Logger {
	return f.logger.zapLogger
}

func (f *Factory) dealWithArgs(message string, args ...any) (string, []zap.Field) {
	var (
		fields []zap.Field
		as     []any
	)

	for _, arg := range args {
		if f, ok := arg.(zap.Field); ok {
			fields = append(fields, f)
		} else {
			as = append(as, arg)
		}
	}

	if len(as) > 0 {
		message = fmt.Sprintf(message, as...)
	}

	if len(fields) > 0 {
		fields = append([]zap.Field{zap.Namespace(f.logOptions.fieldNamespace)}, fields...)
	}

	return message, fields
}

func (f *Factory) Bg() Logger {
	var fields []zap.Field

	opts := f.logOptions
	if opts.addHandler != nil {
		fields = append(fields, opts.addHandler()...)
	}

	return f.logger.With(fields...)
}

func (f *Factory) With(fields ...zap.Field) *Factory {
	return &Factory{
		logger:     f.logger.With(fields...).(*logger),
		logOptions: f.logOptions,
	}
}

func (f *Factory) SetLevel(level string) {
	if f.logger.atomicLevel == nil {
		return
	}

	if err := f.logger.atomicLevel.UnmarshalText([]byte(level)); err != nil {
		Warn("set log level failed", zap.Error(err), zap.String("level", level))
	}
}

func (f *Factory) Err(err error) *Factory {
	if err == nil {
		return f
	}

	var fields []zap.Field

	opts := f.logOptions
	if opts.errHandler != nil {
		fields = append(fields, opts.errHandler(err)...)
	}

	return f.With(fields...)
}

func (f *Factory) Ctx(ctx context.Context) Logger {
	var fields []zap.Field

	opts := f.logOptions
	if opts.ctxHandler != nil {
		fields = append(fields, opts.ctxHandler(ctx)...)
	}

	if len(fields) > 0 {
		return f.Bg().With(fields...)
	}

	return f.Bg()
}
