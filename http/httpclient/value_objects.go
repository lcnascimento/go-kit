package httpclient

import (
	"net/http"

	"github.com/lcnascimento/go-kit/errors"
)

// Headers is a map containing the relation key=value of the headers used on the http rest request.
type Headers map[string]string

// PathParams is a map containing the relation key=value of the path params used on the http rest request.
// It will be used to replace values given in Path parameter.
type PathParams map[string]string

// QueryParams is a map containing the relation key=value of the query params used on the http rest request.
type QueryParams map[string]string

// Request are the params used to build a new http rest request.
type Request struct {
	Host        string
	Path        string
	Body        []byte
	Headers     Headers
	QueryParams QueryParams
	PathParams  PathParams

	acceptStatusCodes map[int]bool
}

// Result are the params returned from the client HTTP request.
type Result struct {
	StatusCode int
	Response   []byte
}

func statusCodeToKind(code int) errors.KindType {
	switch code {
	case http.StatusBadRequest:
		return errors.KindInvalidInput
	case http.StatusUnauthorized:
		return errors.KindUnauthenticated
	case http.StatusForbidden:
		return errors.KindUnauthorized
	case http.StatusNotFound:
		return errors.KindNotFound
	case http.StatusConflict:
		return errors.KindConflict
	case http.StatusInternalServerError:
		return errors.KindInternal
	case http.StatusTooManyRequests:
		return errors.KindResourceExhausted
	case http.StatusServiceUnavailable:
		return errors.KindServiceUnavailable
	case http.StatusUnprocessableEntity:
		return errors.KindUnprocessable
	default:
		return errors.KindInternal
	}
}
