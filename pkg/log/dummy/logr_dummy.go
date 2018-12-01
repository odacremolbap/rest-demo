package dummy

import "github.com/go-logr/logr"

// Logger is a logr that does nothing.
// It will be returned at loggers with verbosity level above configured
type Logger struct{}

var _ logr.Logger = &Logger{}

func (l *Logger) Enabled() bool                                 { return false }
func (l *Logger) Info(_ string, _ ...interface{})               {}
func (l Logger) Error(err error, msg string, kv ...interface{}) {}
func (l *Logger) V(level int) logr.InfoLogger                   { return l }
func (l *Logger) WithName(name string) logr.Logger              { return l }
func (l *Logger) WithValues(kv ...interface{}) logr.Logger      { return l }
