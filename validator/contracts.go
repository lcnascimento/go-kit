package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// CustomValidator defines the interface that custom validators must implement.
// It requires methods to return the validation tag, function, and translation details.
type CustomValidator interface {
	// Tag returns the tag identifier used in struct field validation tags (e.g., `validate:"tag"`).
	Tag() string

	// Func returns the validator.Func that performs the validation logic.
	Func() validator.Func
}

// Validator validates structs given validate annotations.
// Based on https://github.com/go-playground/validator.
type Validator interface {
	// Validate validates a struct given validate annotations.
	Validate(ctx context.Context, s any) error
}
