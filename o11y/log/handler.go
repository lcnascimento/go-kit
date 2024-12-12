package log

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewHandler creates a new slog.Handler with OpenTelemetry support.
func NewHandler(pkg string) slog.Handler {
	return &handler{
		level:   getLevel(),
		Handler: otelslog.NewHandler(pkg),
	}
}

type handler struct {
	*otelslog.Handler
	level slog.Level
}

// Enabled returns true if the handler is enabled for the given level.
func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func getLevel() slog.Level {
	name := "INFO"
	if l := os.Getenv("LOG_LEVEL"); l != "" {
		name = l
	}

	return levelByName[strings.ToUpper(name)]
}
