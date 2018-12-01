package log

import "github.com/go-logr/logr"

// Placeholder for root logger that can be referenced from
// all server code.

// logger hosts the root logger
// Should be initialized at server start before any logging
var logger logr.Logger

// SetDefaultLogger sets the default logger
func SetDefaultLogger(l logr.Logger) {
	logger = l
}

// Error logs an error
func Error(err error, msg string, kv ...interface{}) {
	logger.Error(err, msg, kv...)
}

// Info logs a message
func Info(msg string, kv ...interface{}) {
	logger.Info(msg, kv...)
}

// V return an InfoLogger at the verbosity level
func V(level int) logr.InfoLogger {
	return logger.V(level)
}

// WithName return an InfoLogger with a prefix
func WithName(name string) logr.Logger {
	return logger.WithName(name)
}

// WithValues return an InfoLogger with key value tuples
func WithValues(kv ...interface{}) logr.Logger {
	return logger.WithValues(kv...)
}
