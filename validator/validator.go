package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type val struct {
	validator *validator.Validate
}

// New creates a new validator.
func New(options ...Option) Validator {
	v := &val{
		validator: validator.New(),
	}

	for _, option := range options {
		option(v)
	}

	return v
}

// Validate validates a struct given validate annotations.
func (v *val) Validate(ctx context.Context, s any) error {
	if err := v.validator.StructCtx(ctx, s); err != nil {
		return onInvalidStruct(ctx, s, err)
	}

	return nil
}
