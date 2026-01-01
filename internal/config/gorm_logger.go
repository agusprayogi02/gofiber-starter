package config

import (
	"context"
	"fmt"
	"time"

	"starter-gofiber/pkg/logger"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

// GormLogger implements gorm's logger.Interface with zap
type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormLogger.LogLevel
}

// NewGormLogger creates a new GORM logger
func NewGormLogger(slowThreshold time.Duration, logLevel gormLogger.LogLevel) *GormLogger {
	return &GormLogger{
		SlowThreshold: slowThreshold,
		LogLevel:      logLevel,
	}
}

// LogMode sets log level
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace logs SQL queries
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("duration", elapsed),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error:
		fields = append(fields, zap.Error(err))
		logger.Error("Database Query Error", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		logger.Warn("Slow Database Query", fields...)
	case l.LogLevel >= gormLogger.Info:
		logger.Debug("Database Query", fields...)
	}
}
