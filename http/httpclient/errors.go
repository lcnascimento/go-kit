package httpclient

import (
	"github.com/lcnascimento/go-kit/errors"
)

var ErrUnexpectedStatusCode = errors.New("unexpected status code").WithCode("UNEXPECTED_STATUS_CODE")
