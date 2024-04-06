package logx

import (
	"context"
	"go.uber.org/zap"
)

func writeLog(logMethod func(f *Factory) func(string, ...zap.Field), format string, args ...interface{}) {
	f := GetGlobalFactory()
	message, fields := f.dealWithArgs(format, args...)
	logMethod(f)(message, fields...)
}

func Debug(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Debug
	}, format, args...)
}

func Info(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Info
	}, format, args...)
}

func Warn(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Warn
	}, format, args...)
}

func Error(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Error
	}, format, args...)
}

func Fatal(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Fatal
	}, format, args...)
}

func Panic(format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Bg().Panic
	}, format, args...)
}

func InfoCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Info
	}, format, args...)
}

func DebugCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Debug
	}, format, args...)
}

func WarnCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Warn
	}, format, args...)
}

func ErrorCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Error
	}, format, args...)
}

func FatalCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Fatal
	}, format, args...)
}

func PanicCtx(ctx context.Context, format string, args ...interface{}) {
	writeLog(func(f *Factory) func(string, ...zap.Field) {
		return f.Ctx(ctx).Panic
	}, format, args...)
}
