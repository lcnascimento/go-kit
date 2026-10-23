package errortracking

import (
	"fmt"
	"log/slog"
)

type warnCappedLogger struct{}

func (warnCappedLogger) Debugf(format string, args ...any) {
	slog.Default().Debug(fmt.Sprintf(format, args...))
}

func (warnCappedLogger) Logf(format string, args ...any) {
	slog.Default().Warn(fmt.Sprintf(format, args...))
}

func (warnCappedLogger) Warnf(format string, args ...any) {
	slog.Default().Warn(fmt.Sprintf(format, args...))
}

func (warnCappedLogger) Errorf(format string, args ...any) {
	slog.Default().Warn(fmt.Sprintf(format, args...))
}
