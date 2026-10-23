package errortracking

import (
	"fmt"
	"log/slog"
)

// warnLogger never logs above Warn, so the SDK's own failures are never captured as exceptions.
type warnLogger struct{}

func (warnLogger) Debugf(format string, args ...any) {
	slog.Default().Debug(fmt.Sprintf(format, args...))
}

func (warnLogger) Logf(format string, args ...any) {
	slog.Default().Info(fmt.Sprintf(format, args...))
}

func (warnLogger) Warnf(format string, args ...any) {
	slog.Default().Warn(fmt.Sprintf(format, args...))
}

func (warnLogger) Errorf(format string, args ...any) {
	slog.Default().Warn(fmt.Sprintf(format, args...))
}
