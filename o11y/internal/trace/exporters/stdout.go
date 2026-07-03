package exporters

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

// Stdout creates a new stdout trace exporter.
func Stdout() (trace.SpanExporter, error) {
	var opts []stdouttrace.Option

	cfg := global.Config()

	if cfg.PrettyPrint {
		opts = append(opts, stdouttrace.WithPrettyPrint())
	}

	exporter, err := stdouttrace.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stdout trace exporter")
	}

	return exporter, nil
}
