package logger

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// InitSentry initializes Sentry for error tracking
func InitSentry(dsn, env string) error {
	if dsn == "" {
		return nil // Sentry is optional
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      env,
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return event
		},
	})

	if err != nil {
		return err
	}

	Info("Sentry initialized", zap.String("env", env))
	return nil
}

// FlushSentry flushes any pending Sentry events
func FlushSentry() {
	sentry.Flush(time.Second * 2)
}

// CaptureException captures an exception and sends it to Sentry
func CaptureException(err error) {
	if err != nil {
		sentry.CaptureException(err)
	}
}

// CaptureMessage captures a message and sends it to Sentry
func CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}
