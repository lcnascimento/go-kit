package exporters

import (
	"encoding/json"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

// Stdout creates a new stdout metric reader.
func Stdout() (metric.Exporter, error) {
	cfg := global.Config()

	enc := json.NewEncoder(os.Stdout)
	opts := []stdoutmetric.Option{
		stdoutmetric.WithEncoder(enc),
		stdoutmetric.WithoutTimestamps(),
	}

	if cfg.PrettyPrint {
		enc.SetIndent("", "  ")
	}

	exporter, err := stdoutmetric.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stdout metric exporter")
	}

	return exporter, nil
}
