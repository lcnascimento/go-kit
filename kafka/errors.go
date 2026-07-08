package kafka

import "github.com/lcnascimento/go-kit/errors"

var ErrWriteMessages = errors.New("failed to write message(s) to event stream").
	WithCode("ERR_WRITE_MESSAGES").
	Retryable()
