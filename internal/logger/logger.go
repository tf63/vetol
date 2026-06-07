package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	log = slog.New(slog.NewTextHandler(os.Stderr, nil))
}

// Error logs an error-level message.
func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

// Warn logs a warning-level message.
func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

// Info logs an info-level message.
func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

// Debug logs a debug-level message.
func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}
