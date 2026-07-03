package profile

import (
	"context"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

var prof Profiler

// Profiler is an interface for profiling the application.
type Profiler interface {
	// Start starts the profiler.
	Start(ctx context.Context) error

	// Shutdown shuts down the profiler.
	Shutdown(ctx context.Context) error
}

// Start starts the profiler.
func Start(ctx context.Context) error {
	cfg := global.Config()

	if cfg.Disabled || cfg.ProfilesExporter == config.ProfilesExporterNone {
		return nil
	}

	return errors.New("unsupported profiler: %s", cfg.ProfilesExporter)
}

// Shutdown shuts down the profiler.
func Shutdown(ctx context.Context) error {
	if prof == nil {
		return nil
	}

	return prof.Shutdown(ctx)
}
