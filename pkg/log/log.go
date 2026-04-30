package log

import (
	"fmt"

	flogger "github.com/mirkobrombin/go-foundation/pkg/logger"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelSuccess
	LogLevelWarning
	LogLevelError
)

// Logger defines the interface for logging messages with different severity levels.
type Logger interface {
	Info(format string, a ...any)
	Success(format string, a ...any)
	Warning(format string, a ...any)
	Error(format string, a ...any)
}

type wrappedLogger struct {
	inner flogger.Logger
}

// New creates a new instance of the Logger, backed by go-foundation's structured logger.
//
// Example:
//
//	logger := log.New()
//	logger.Info("Starting application...")
func New() Logger {
	return &wrappedLogger{inner: flogger.New()}
}

func (l *wrappedLogger) Info(format string, a ...any) {
	l.inner.Info(fmt.Sprintf(format, a...))
}

func (l *wrappedLogger) Success(format string, a ...any) {
	l.inner.Info(fmt.Sprintf(format, a...))
}

func (l *wrappedLogger) Warning(format string, a ...any) {
	l.inner.Warn(fmt.Sprintf(format, a...))
}

func (l *wrappedLogger) Error(format string, a ...any) {
	l.inner.Error(fmt.Sprintf(format, a...))
}

// NewFoundationLogger returns the underlying go-foundation structured logger,
// useful for advanced use cases like adding sinks or structured fields.
//
// Example:
//
//	fl := log.NewFoundationLogger()
//	fl.Info("structured message", logger.Field{Key: "user", Value: "alice"})
func NewFoundationLogger() flogger.Logger {
	return flogger.New()
}