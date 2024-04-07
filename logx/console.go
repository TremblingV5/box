package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var globalConsoleLogger = getConsoleLogger()

func getConsoleLogger() *zap.Logger {
	configs := zap.NewDevelopmentConfig()
	zapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(configs.EncoderConfig),
		zapcore.AddSync(os.Stderr),
		zap.DebugLevel,
	)

	return zap.New(zapCore, zap.AddCallerSkip(3), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func Console() *zap.Logger {
	return globalConsoleLogger
}
