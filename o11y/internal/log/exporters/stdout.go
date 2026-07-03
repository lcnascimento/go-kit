package exporters

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

// Stdout creates a new log stdout exporter.
func Stdout() (log.Exporter, error) {
	cfg := global.Config()

	var opts []stdoutlog.Option
	if cfg.PrettyPrint {
		opts = append(opts, stdoutlog.WithPrettyPrint())
	}

	exporter, err := stdoutlog.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stdout log exporter")
	}

	return exporter, nil
}
