package log

import (
	"time"

	"github.com/lcnascimento/go-kit/propagation"
)

// Level indicates the severity of the data being logged.
type Level int

const (
	// LevelCritical alerts about severe problems. Most of the time, needs some human intervention ASAP.
	LevelCritical Level = iota + 1
	// LevelError alerts about events that are likely to cause problems.
	LevelError
	// LevelWarning warns about events the might cause problems to the system.
	LevelWarning
	// LevelInfo are routine information.
	LevelInfo
	// LevelDebug are debug or trace information.
	LevelDebug
)

var levelStringValueMap = map[string]Level{
	"CRITICAL": LevelCritical,
	"ERROR":    LevelError,
	"WARNING":  LevelWarning,
	"INFO":     LevelInfo,
	"DEBUG":    LevelDebug,
}

// String returns the name of the LogLevel.
func (l Level) String() string {
	return []string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"INFO",
		"DEBUG",
	}[l-1]
}

// Option is a type to set Logger options.
type Option func(*Logger)

// WithLevel instructs the Logger to use the given log level.
func WithLevel(level string) Option {
	return func(l *Logger) {
		l.level = levelStringValueMap[level]
	}
}

// WithContextKeySet instructs the Logger to include the given context keys in logs.
func WithContextKeySet(keys propagation.ContextKeySet) Option {
	return func(l *Logger) {
		l.contextKeys = keys
	}
}

// WithFormatter instructs the Logger to use the given formatter for logging.
func WithFormatter(formatter Formatter) Option {
	return func(l *Logger) {
		l.formatter = formatter
	}
}

// WithTimmer instructs the Logger to use the given timmer to get current Time for logging.
// It must be used just for testing.
func WithTimmer(now func() time.Time) Option {
	return func(l *Logger) {
		l.now = now
	}
}
