package log

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewHandler creates a new slog.Handler with OpenTelemetry support.
func NewHandler(pkg string) slog.Handler {
	return otelslog.NewHandler(pkg)
}
