package logrus

import (
	"fmt"

	lgrus "github.com/Sirupsen/logrus"
	"github.com/go-logr/logr"
	"github.com/odacremolbap/rest-demo/pkg/log/dummy"
)

const defaultVerbosity = 1

// Logger encapsulates logrus in logr
type Logger struct {
	name      string
	log       *lgrus.Logger
	entry     *lgrus.Entry
	verbosity int
}

var notActivated = &dummy.Logger{}
var _ logr.Logger = &Logger{}
var logVerbosity int

// // NewLogrusStandardLogr returns an encapsulated standard logrus logr
// func NewLogrusStandardLogr(verbosity int) *Logger {
// 	return NewLogrusCustomLogr(lgrus.StandardLogger(), "", verbosity)
// }

// NewLogrusLogr returns the logrus logger encapsulated into a logr interface
func NewLogrusLogr(l *lgrus.Logger, name string, verbosity int) *Logger {
	logVerbosity = verbosity
	return &Logger{
		log:       l,
		name:      name,
		verbosity: defaultVerbosity,
		entry:     nil,
	}
}

// Enabled checks if the logger is enabled upon its verbosity
func (l Logger) Enabled() bool {
	return l.verbosity <= logVerbosity
}

func (l Logger) Error(err error, msg string, kv ...interface{}) {
	kv = append(kv, "level", "error")
	if len(l.name) != 0 {
		kv = append(kv, "logger", l.name)
	}
	entry := l.entry

	newFields := l.fieldsFromKV(kv...)
	if entry != nil {
		entry = entry.WithFields(*newFields)
	} else {
		entry = l.log.WithFields(*newFields)
	}

	if err != nil {
		if entry != nil {
			entry = entry.WithError(err)
		} else {
			entry = l.log.WithError(err)
		}
	}

	if entry != nil {
		entry.Error(msg)
		return
	}
	l.log.WithError(err).Error(msg)
}

// Info logs non error messages
func (l *Logger) Info(msg string, kv ...interface{}) {
	if len(l.name) != 0 {
		kv = append(kv, "logger", l.name)
	}
	if len(kv) != 0 {
		newFields := l.fieldsFromKV(kv...)
		if l.entry != nil {
			l.entry.WithFields(*newFields).Info(msg)
			return
		}
		l.log.WithFields(*newFields).Info(msg)
		return
	}

	if l.entry != nil {
		l.entry.Info(msg)
		return
	}
	l.log.Info(msg)
}

// V return logger with verbosity level set
func (l *Logger) V(level int) logr.InfoLogger {
	if level > logVerbosity {
		return notActivated
	}
	return &Logger{
		name:      l.name,
		log:       l.log,
		entry:     l.entry,
		verbosity: level,
	}
}

// WithName returns a logger with a name prefix
func (l *Logger) WithName(name string) logr.Logger {
	if len(l.name) > 0 {
		name = l.name + "." + name
	}
	newLog := l.copyLogger()
	newLog.name = name
	return newLog
}

func (l *Logger) copyLogger() *Logger {

	return &Logger{
		name:      l.name,
		log:       l.log,
		entry:     l.entry,
		verbosity: l.verbosity,
	}
}

func (l *Logger) fieldsFromKV(kv ...interface{}) *lgrus.Fields {
	newFields := lgrus.Fields{}
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			panic(fmt.Sprintf("key writen at KV pair to log is not a string: %+v\n", kv[i]))
		}
		if i+1 >= len(kv) {
			// avoid out of range
			newFields[k] = ""
			break
		}
		newFields[k] = kv[i+1]
	}
	return &newFields
}

// WithValues adds key/value pairs to the log entry
func (l *Logger) WithValues(kv ...interface{}) logr.Logger {
	if len(kv) == 0 {
		return l
	}
	newFields := l.fieldsFromKV(kv)
	newLogger := l.copyLogger()

	if newLogger.entry != nil {
		newLogger.entry = newLogger.entry.WithFields(*newFields)
	} else {
		newLogger.entry = newLogger.log.WithFields(*newFields)
	}
	return newLogger
}
