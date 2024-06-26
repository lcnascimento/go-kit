package httpclient

import (
	"net/http"

	"github.com/lcnascimento/go-kit/errors"
)

// All errors that can be emitted by http requests.
var (
	ErrParseURL = func(err error) error {
		return errors.New("error on parsing the request url").
			WithCode("ERR_PARSE_REQUEST_URL").
			WithKind(errors.KindInvalidInput).
			WithCause(err)
	}

	ErrBuildRequestError = func(err error) error {
		return errors.New("could not build http request").
			WithCode("ERR_BUILD_HTTP_REQUEST").
			WithKind(errors.KindInternal).
			WithCause(err)
	}

	ErrRequestError = func(err error) error {
		return errors.New("http request error").
			WithCode("HTTP_REQUEST_ERROR").
			WithKind(errors.KindUnknown).
			WithCause(err).
			Retryable()
	}

	ErrBodyReadError = func(err error) error {
		return errors.New("http request error").
			WithCode("ERR_READ_HTTP_RESPONSE_BODY").
			WithKind(errors.KindUnknown).
			WithCause(err)
	}

	ErrUnexpectedStatusCode = func(code int) error {
		kind := kindByStatusCode(code)
		retry := retryByStatusCode[code]

		err := errors.New("unexpected http response status code").
			WithCode("ERR_UNEXPECTED_STATUS_CODE").
			WithKind(kind)

		if !retry {
			return err
		}

		return err.Retryable()
	}

	ErrBodyCastError = func(err error) error {
		return errors.New("could not cast http response body").
			WithCode("ERR_HTTP_RESPONSE_BODY_CAST").
			WithKind(errors.KindUnknown).
			WithCause(err)
	}
)

func kindByStatusCode(code int) errors.KindType {
	switch code {
	case http.StatusBadRequest:
		return errors.KindInvalidInput
	case http.StatusUnauthorized:
		return errors.KindUnauthenticated
	case http.StatusForbidden:
		return errors.KindUnauthorized
	case http.StatusNotFound:
		return errors.KindNotFound
	case http.StatusTooManyRequests:
		return errors.KindResourceExhausted
	default:
		return errors.KindUnknown
	}
}

var retryByStatusCode = map[int]bool{
	http.StatusInternalServerError: true,
	http.StatusBadGateway:          true,
	http.StatusServiceUnavailable:  true,
	http.StatusGatewayTimeout:      true,
}
