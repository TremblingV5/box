package mysqlx

import (
	"context"
	"errors"
	"fmt"
	"github.com/TremblingV5/box/internal/nlog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type loggerAdaptor struct {
	log nlog.Logger
	logger.Config
}

func newLoggerAdaptor() logger.Interface {
	return &loggerAdaptor{
		log: nlog.GetOrNewLogger("mysql-gorm"),
		Config: logger.Config{
			SlowThreshold: 0,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	}
}

func (l *loggerAdaptor) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *loggerAdaptor) Info(ctx context.Context, s string, i ...interface{}) {
	l.log.LogKVDepth(1, nlog.InfoLevel, fmt.Sprintf(s, i...), nlog.KVCtx(ctx))
}

func (l *loggerAdaptor) Warn(ctx context.Context, s string, i ...interface{}) {
	l.log.LogKVDepth(1, nlog.WarnLevel, fmt.Sprintf(s, i...), nlog.KVCtx(ctx))
}

func (l *loggerAdaptor) Error(ctx context.Context, s string, i ...interface{}) {
	l.log.LogKVDepth(1, nlog.ErrorLevel, fmt.Sprintf(s, i...), nlog.KVCtx(ctx))
}

func (l *loggerAdaptor) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	var (
		nlogLevel nlog.Level
		nlogMsg   = "gorm log"
		elapsed   = time.Since(begin)
		sql, rows = fc()
		rowsVal   any
	)

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			nlogLevel = nlog.ErrorLevel
		} else if l.LogLevel == logger.Info {
			nlogLevel = nlog.ErrorLevel
		} else {
			return
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		nlogLevel = nlog.WarnLevel
		nlogMsg = nlogMsg + fmt.Sprintf(", SLOW SQL >= %v", l.SlowThreshold)
	case l.LogLevel == logger.Info:
		nlogLevel = nlog.InfoLevel
	default:
		return
	}

	if rows == -1 {
		rowsVal = "-"
	} else {
		rowsVal = rows
	}

	l.log.LogKVDepth(
		1,
		nlogLevel,
		nlogMsg,
		nlog.KVCtx(ctx),
		nlog.KVVAny("rows", rowsVal),
		nlog.KVVAny("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
		nlog.KVString("sql", sql),
		nlog.KVError(err),
	)
}
