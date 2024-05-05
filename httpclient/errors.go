package httpclient

import "github.com/lcnascimento/go-kit/errors"

// ErrMetricInitialization occurs when a metric initialization fails.
var ErrMetricInitialization = func(metric string) error {
	return errors.New("could not initialize metric: %s", metric).
		WithCode("ERR_METRIC_INITIALIZATION").
		WithKind(errors.KindUnexpected).
		WithStack().
		Retryable(false)
}
