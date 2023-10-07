package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Log *logrus.Logger
}

func NewGormLogger(log *logrus.Logger) *GormLogger {
	return &GormLogger{Log: log}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Log.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Log.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Log.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Log.IsLevelEnabled(logrus.TraceLevel) {
		elapsed := time.Since(begin)
		sql, rows := fc()
		fields := logrus.Fields{
			"sql":        sql,
			"rows":       rows,
			"elapsed_ms": float64(elapsed.Nanoseconds()) / 1e6,
		}
		if err != nil {
			l.Log.WithContext(ctx).WithFields(fields).WithError(err).Trace("gorm: error")
		} else {
			l.Log.WithContext(ctx).WithFields(fields).Trace("gorm: trace")
		}
	}
}
