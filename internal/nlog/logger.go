package nlog

import (
	"context"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	TraceLevel Level = zapcore.DebugLevel - 1
	DebugLevel Level = zapcore.DebugLevel
	InfoLevel  Level = zapcore.InfoLevel
	WarnLevel  Level = zapcore.WarnLevel
	ErrorLevel Level = zapcore.ErrorLevel
	PanicLevel Level = zapcore.PanicLevel
	FatalLevel Level = zapcore.FatalLevel
)

type Logger interface {
	BaseLogger
	LevelLogger
	WithKV(kv ...KV) Logger
}

type LevelLogger interface {
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)
	TraceKV(message string, kv ...KV)
	DebugKV(message string, kv ...KV)
	InfoKV(message string, kv ...KV)
	WarnKV(message string, kv ...KV)
	ErrorKV(message string, kv ...KV)
	PanicKV(message string, kv ...KV)
	FatalKV(message string, kv ...KV)
}

type BaseLogger interface {
	Name() string
	Enabled(level Level) bool
	Logf(level Level, format string, args ...any)
	LogKV(level Level, message string, kv ...KV)
	LogfDepth(callDepth int, level Level, format string, args ...any)
	LogKVDepth(callDepth int, level Level, message string, kv ...KV)
}

type Kind int

const (
	KindAny Kind = iota
	KindError
	KindContext
)

type KV struct {
	Key   string
	Kind  Kind
	Value any
}

func KVVAny(key string, value any) KV {
	return KV{
		Key:   key,
		Kind:  KindAny,
		Value: value,
	}
}

func KVString(key, value string) KV {
	return KV{
		Key:   key,
		Kind:  KindAny,
		Value: value,
	}
}

type Numeric interface {
	int | int32 | int64 | uint | uint32 | uint64 | float32 | float64 | int8 | int16 | uint8 | uint16
}

func KVNumber[T Numeric](key string, value T) KV {
	return KV{
		Key:   key,
		Kind:  KindAny,
		Value: value,
	}
}

func KVBool(key string, value bool) KV {
	return KV{
		Key:   key,
		Kind:  KindAny,
		Value: value,
	}
}

func KVCtx(ctx context.Context) KV {
	return KV{
		Key:   "context",
		Kind:  KindContext,
		Value: ctx,
	}
}

func KVError(err error) KV {
	return KV{
		Key:   "error",
		Kind:  KindError,
		Value: err,
	}
}
