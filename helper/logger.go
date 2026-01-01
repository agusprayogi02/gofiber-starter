package helper

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global logger instance
// This logger is used both for application logging and Fiber HTTP logging via fiberzap
func InitLogger(env string) error {
	var config zap.Config

	if env == "prod" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
	}

	// Set output paths
	config.OutputPaths = []string{"stdout"}
	if env == "prod" {
		config.OutputPaths = append(config.OutputPaths, "logs/app.log")
	}

	// Create logger
	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Sync flushes any buffered log entries
func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// LogRequest logs HTTP request information
func LogRequest(method, path string, statusCode int, duration time.Duration, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
	}

	allFields := append(baseFields, fields...)

	if statusCode >= 500 {
		Error("HTTP Request", allFields...)
	} else if statusCode >= 400 {
		Warn("HTTP Request", allFields...)
	} else {
		Info("HTTP Request", allFields...)
	}
}

// LogError logs error with context
func LogError(err error, context string, fields ...zap.Field) {
	if err == nil {
		return
	}

	baseFields := []zap.Field{
		zap.Error(err),
		zap.String("context", context),
	}

	allFields := append(baseFields, fields...)
	Error("Error occurred", allFields...)
}

// LogDBQuery logs database query information
func LogDBQuery(query string, duration time.Duration, rowsAffected int64, err error) {
	fields := []zap.Field{
		zap.String("query", query),
		zap.Duration("duration", duration),
		zap.Int64("rows_affected", rowsAffected),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Error("Database Query Failed", fields...)
	} else if duration > 1*time.Second {
		Warn("Slow Database Query", fields...)
	} else {
		Debug("Database Query", fields...)
	}
}
