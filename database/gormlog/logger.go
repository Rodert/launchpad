package gormlog

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"launchpad/constant"
)

var _ logger.Interface = &gormLogger{}

var (
	slowSqlStr = ""
)

type gormLogger struct {
	l *logrus.Logger
	logger.Config
}

func (l *gormLogger) getLogger(ctx context.Context, fileLine, sqlStr string, elapsed, rows int64) *logrus.Entry {
	requestId := ctx.Value(constant.CONTEXT_KEY_REQUEST_ID)
	entry := logrus.NewEntry(l.l)
	fields := make(logrus.Fields, 6)
	fields["subModule"] = "gorm"
	fields["fileLine"] = fileLine
	if requestId != nil {
		fields[constant.CONTEXT_KEY_REQUEST_ID] = requestId
	}
	if len(sqlStr) > 0 {
		fields["sqlStr"] = sqlStr
	}
	if elapsed > 0 {
		fields["elapsed"] = fmt.Sprintf("%.3fms", float64(elapsed)/1e6)
	}
	if rows != -1 {
		fields["rows"] = rows
	}
	entry = entry.WithFields(fields)
	return entry
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log := l.getLogger(ctx, utils.FileWithLineNum(), "", 0, -1)
		log.Infof(msg, data...)
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log := l.getLogger(ctx, utils.FileWithLineNum(), "", 0, -1)
		log.Warnf(msg, data...)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log := l.getLogger(ctx, utils.FileWithLineNum(), "", 0, -1)
		log.Errorf(msg, data...)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		elapsed := time.Since(begin)
		sql, rows := fc()
		if elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn {
			sql += slowSqlStr
		}
		log := l.getLogger(ctx, utils.FileWithLineNum(), sql, elapsed.Nanoseconds(), rows)
		log.Info("sql trace")
	}
}

func New(l *logrus.Logger, config logger.Config) logger.Interface {
	slowSqlStr = fmt.Sprintf("[SLOW SQL:%v]", config.SlowThreshold)
	return &gormLogger{
		l:      l,
		Config: config,
	}
}
