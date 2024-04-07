package nlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"reflect"
)

func init() {
	SetLoggerBuilder(func(name string) Logger {
		return NewLoggerFromZap(NewDevZapLogger(zapcore.DebugLevel), name)
	})
}

var _ Logger = (*logger)(nil)

type logger struct {
	name       string
	optionList *optionList
	zapLogger  *zap.Logger
	kvs        []KV
}

func (s *logger) Name() string {
	return s.name
}

func (s *logger) Enabled(level Level) bool {
	return s.zapLogger.Core().Enabled(level)
}

const (
	defaultNamespace = "content"
)

func (s *logger) processKVs(kvs ...KV) (fields []zap.Field) {
	if s.optionList.rootKVAdder != nil {
		for _, kv := range s.optionList.rootKVAdder() {
			fields = append(fields, zap.Any(kv.Key, kv.Value))
		}
	}

	namespace := []zap.Field{
		zap.Namespace(defaultNamespace),
	}

	kvs = append(s.kvs, kvs...)

	for _, v := range kvs {
		switch v.Kind {
		case KindAny:
			fields = append(fields, zap.Any(v.Key, v.Value))
		default:
			if s.optionList.kvConverter != nil {
				for _, kv := range s.optionList.kvConverter(v) {
					fields = append(fields, zap.Any(kv.Key, kv.Value))
				}
			}
		}
	}

	if len(namespace) > 1 {
		fields = append(fields, namespace...)
	}

	return fields
}

func (s *logger) Logf(level Level, format string, args ...any) {
	if !s.Enabled(level) {
		return
	}

	s.zapLogger.Log(level, fmt.Sprintf(format, args...), s.processKVs()...)
}

func (s *logger) LogKV(level Level, message string, kv ...KV) {
	if !s.Enabled(level) {
		return
	}

	s.zapLogger.Log(level, message, s.processKVs(kv...)...)
}

func (s *logger) LogfDepth(callDepth int, level Level, format string, args ...any) {
	if !s.Enabled(level) {
		return
	}

	s.zapLogger.WithOptions(zap.AddCallerSkip(callDepth)).Log(level, fmt.Sprintf(format, args...), s.processKVs()...)
}

func (s *logger) LogKVDepth(callDepth int, level Level, message string, kv ...KV) {
	if !s.Enabled(level) {
		return
	}

	s.zapLogger.WithOptions(zap.AddCallerSkip(callDepth)).Log(level, message, s.processKVs(kv...)...)
}

func (s *logger) clone() *logger {
	c := *s
	return &c
}

func (s *logger) WithKV(kv ...KV) Logger {
	c := s.clone()
	c.kvs = append(c.kvs, kv...)
	return c
}

func (s *logger) Tracef(format string, args ...any) {
	s.Logf(TraceLevel, format, args...)
}

func (s *logger) Debugf(format string, args ...any) {
	s.Logf(DebugLevel, format, args...)
}

func (s *logger) Infof(format string, args ...any) {
	s.Logf(InfoLevel, format, args...)
}

func (s *logger) Warnf(format string, args ...any) {
	s.Logf(WarnLevel, format, args...)
}

func (s *logger) Errorf(format string, args ...any) {
	s.Logf(ErrorLevel, format, args...)
}

func (s *logger) Panicf(format string, args ...any) {
	s.Logf(PanicLevel, format, args...)
}

func (s *logger) Fatalf(format string, args ...any) {
	s.Logf(FatalLevel, format, args...)
}

func (s *logger) TraceKV(message string, kv ...KV) {
	s.LogKV(TraceLevel, message, kv...)
}

func (s *logger) DebugKV(message string, kv ...KV) {
	s.LogKV(DebugLevel, message, kv...)
}

func (s *logger) InfoKV(message string, kv ...KV) {
	s.LogKV(InfoLevel, message, kv...)
}

func (s *logger) WarnKV(message string, kv ...KV) {
	s.LogKV(WarnLevel, message, kv...)
}

func (s *logger) ErrorKV(message string, kv ...KV) {
	s.LogKV(ErrorLevel, message, kv...)
}

func (s *logger) PanicKV(message string, kv ...KV) {
	s.LogKV(PanicLevel, message, kv...)
}

func (s *logger) FatalKV(message string, kv ...KV) {
	s.LogKV(FatalLevel, message, kv...)
}

func NewLoggerFromZap(l *zap.Logger, name string, options ...Option) Logger {
	optionList := defaultOptionList()
	for _, option := range options {
		option(optionList)
	}

	return &logger{
		name:       name,
		optionList: optionList,
		zapLogger:  l.Named(name),
	}
}

func newZapLoggerFromCore(zapCore zapcore.Core) *zap.Logger {
	return zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(3), zap.AddStacktrace(zapcore.ErrorLevel))
}

func NewDevZapLogger(level zapcore.LevelEnabler) *zap.Logger {
	encConfig := zap.NewDevelopmentEncoderConfig()
	encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encConfig.FunctionKey = "F"
	zapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encConfig),
		zapcore.AddSync(os.Stderr),
		level,
	)

	return newZapLoggerFromCore(zapCore)
}

func NewProdZapLogger(level zapcore.LevelEnabler, wr, errWr io.Writer) *zap.Logger {
	encConfig := zap.NewProductionEncoderConfig()
	encConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encConfig.FunctionKey = "func"

	var cores []zapcore.Core

	if errWr == nil || reflect.ValueOf(errWr).IsNil() {
		defaultCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encConfig),
			zapcore.AddSync(errWr),
			level,
		)
		cores = append(cores, defaultCore)
	} else {
		errorCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encConfig),
			zapcore.AddSync(wr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zap.ErrorLevel && level.Enabled(lvl)
			}),
		)
		cores = append(cores, errorCore)

		defaultCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encConfig),
			zapcore.AddSync(errWr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zap.ErrorLevel && level.Enabled(lvl)
			}),
		)
		cores = append(cores, defaultCore)
	}

	return newZapLoggerFromCore(zapcore.NewTee(cores...))
}
