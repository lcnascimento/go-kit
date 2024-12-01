//nolint:errorlint // ok
package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/lcnascimento/go-kit/errors"
)

// ErrInvalidStruct is the error returned when the struct is invalid.
var ErrInvalidStruct = func(err error) error {
	var out errors.CustomError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			for _, ie := range strings.Split(e.Error(), "\n") {
				out, _ = errors.Wrap(out, ie).(errors.CustomError)
			}
		}
	}

	out, _ = errors.Wrap(out, "invalid struct").(errors.CustomError)

	return out.WithKind(errors.KindInvalidInput)
}
