package profile

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/grafana/pyroscope-go"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/config"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

const mutexProfileFraction = 5

var prof *pyroscope.Profiler

var profileTypes = map[config.Profile][]pyroscope.ProfileType{
	config.ProfileCPU: {pyroscope.ProfileCPU},
	config.ProfileHeap: {
		pyroscope.ProfileAllocObjects,
		pyroscope.ProfileAllocSpace,
		pyroscope.ProfileInuseObjects,
		pyroscope.ProfileInuseSpace,
	},
	config.ProfileGoroutine: {pyroscope.ProfileGoroutines},
	config.ProfileMutex:     {pyroscope.ProfileMutexCount, pyroscope.ProfileMutexDuration},
}

// Start starts the profiler.
func Start(_ context.Context) error {
	cfg := global.Config()

	if cfg.Disabled || cfg.ProfilesExporter == config.ProfilesExporterNone {
		return nil
	}

	if cfg.ProfilesExporter != config.ProfilesExporterPyroscope {
		return errors.New("unsupported profiler: %s", cfg.ProfilesExporter)
	}

	types := []pyroscope.ProfileType{}
	for kind, enabled := range cfg.Profiles {
		if enabled {
			types = append(types, profileTypes[kind]...)
		}
	}

	if cfg.Profiles[config.ProfileMutex] {
		runtime.SetMutexProfileFraction(mutexProfileFraction)
	}

	p, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ServiceName,
		ServerAddress:   cfg.PyroscopeEndpoint,
		Tags:            tags(cfg.ResourceAttributes),
		ProfileTypes:    types,
		UploadRate:      cfg.ProfileExportInterval,
		Logger:          slogLogger{},
	})
	if err != nil {
		return errors.Wrap(err, "failed to start pyroscope profiler")
	}

	prof = p
	return nil
}

// Shutdown shuts down the profiler.
func Shutdown(_ context.Context) error {
	if prof == nil {
		return nil
	}

	if err := prof.Stop(); err != nil {
		return errors.Wrap(err, "failed to shutdown pyroscope profiler")
	}

	prof = nil
	return nil
}

func tags(attrs config.ResourceAttributes) map[string]string {
	output := make(map[string]string, len(attrs))
	for key, value := range attrs {
		if key == semconv.ServiceNameKey {
			continue
		}
		output[strings.ReplaceAll(string(key), ".", "_")] = value
	}
	return output
}

type slogLogger struct{}

func (slogLogger) Infof(format string, args ...any) {
	slog.Default().Debug(fmt.Sprintf(format, args...))
}

func (slogLogger) Debugf(format string, args ...any) {
	slog.Default().Debug(fmt.Sprintf(format, args...))
}

func (slogLogger) Errorf(format string, args ...any) {
	slog.Default().Error(fmt.Sprintf(format, args...))
}
