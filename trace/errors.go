package trace

import "github.com/lcnascimento/go-kit/errors"

// ErrSpanExporterNotSupported indicates that the span exporter used for tracing is not supported.
var ErrSpanExporterNotSupported = errors.New("span exporter not supported").
	WithCode("ERR_SPAN_EXPORTER_NOT_SUPPORTED").
	WithKind(errors.KindInvalidInput).
	Retryable(false)
