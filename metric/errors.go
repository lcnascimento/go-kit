package metric

import "github.com/lcnascimento/go-kit/errors"

// ErrMetricReaderNotSupported indicates that the metric reader is not supported.
var ErrMetricReaderNotSupported = errors.New("metric reader not supported").
	WithCode("ERR_METRIC_READER_NOT_SUPPORTED").
	WithKind(errors.KindInvalidInput).
	Retryable(false)
