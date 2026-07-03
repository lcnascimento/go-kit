package validator

import "github.com/lcnascimento/go-kit/errors"

var (
	ErrUnexpectedValidationError = errors.New("unexpected validation error").
		WithCode("ERR_UNEXPECTED_VALIDATION_ERROR").
		WithKind(errors.KindInternal)
)
