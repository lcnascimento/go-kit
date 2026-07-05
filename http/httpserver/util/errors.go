package util

import (
	"context"
	"net/http"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	ErrParseRequestBody     = errors.New("failed to parse request body").WithCode("ERR_PARSE_REQUEST_BODY").WithKind(errors.KindInvalidInput)
	ErrMissingCorrelationID = errors.New("missing correlation id parameter").WithCode("ERR_MISSING_CORRELATION_ID").WithKind(errors.KindInvalidInput)

	logger = log.MustNewLogger("github.com/lcnascimento/go-kit/http/httpserver/util")
)

func NewAPIError(err error) *APIError {
	out := &APIError{
		Code:      string(errors.Code(err)),
		Message:   err.Error(),
		Retryable: errors.IsRetryable(err),
	}

	reasons := errors.SafeReasons(err)
	if len(reasons) > 0 {
		out.Details = map[string]any{
			"reasons": reasons,
		}
	}

	return out
}

type APIError struct {
	Code      string         `example:"ERR_SOMETHING_WENT_WRONG" json:"code"`
	Message   string         `example:"Something went wrong"     json:"message"`
	Retryable bool           `example:"true"                     json:"retryable"`
	Details   map[string]any `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func WriteError(ctx context.Context, rw http.ResponseWriter, err error) {
	kind := errors.Kind(err)
	status := kindToHTTPStatusCode(kind)

	if errors.Is(err, ErrParseRequestBody) {
		logger.Error(ctx, err)
	}

	WriteResponse(rw, status, NewAPIError(err))
}

func kindToHTTPStatusCode(kind errors.KindType) int {
	switch kind {
	case errors.KindInvalidInput:
		return http.StatusBadRequest
	case errors.KindUnauthenticated:
		return http.StatusUnauthorized
	case errors.KindUnauthorized:
		return http.StatusForbidden
	case errors.KindNotFound:
		return http.StatusNotFound
	case errors.KindConflict:
		return http.StatusConflict
	case errors.KindInternal:
		return http.StatusInternalServerError
	case errors.KindResourceExhausted:
		return http.StatusTooManyRequests
	case errors.KindServiceUnavailable:
		return http.StatusServiceUnavailable
	case errors.KindUnprocessable:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
