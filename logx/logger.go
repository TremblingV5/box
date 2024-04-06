package logx

import "go.uber.org/zap"

type Logger interface {
	Info(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Panic(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

type logger struct {
	zapLogger   *zap.Logger
	atomicLevel *zap.AtomicLevel
}

func (l *logger) Info(message string, fields ...zap.Field) {
	l.zapLogger.Info(message, fields...)
}

func (l *logger) Debug(message string, fields ...zap.Field) {
	l.zapLogger.Debug(message, fields...)
}

func (l *logger) Warn(message string, fields ...zap.Field) {
	l.zapLogger.Warn(message, fields...)
}

func (l *logger) Error(message string, fields ...zap.Field) {
	l.zapLogger.Error(message, fields...)
}

func (l *logger) Panic(message string, fields ...zap.Field) {
	l.zapLogger.Panic(message, fields...)
}

func (l *logger) Fatal(message string, fields ...zap.Field) {
	l.zapLogger.Fatal(message, fields...)
}

func (l *logger) clone() *logger {
	c := *l
	return &c
}

func (l *logger) With(fields ...zap.Field) Logger {
	if len(fields) == 0 {
		return l
	}

	c := l.clone()
	c.zapLogger = l.zapLogger.With(fields...)
	return c
}
