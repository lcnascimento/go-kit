package main

import (
	"context"

	"github.com/lcnascimento/go-kit/validator"
)

func main() {
	ctx := context.Background()
	v := validator.New()

	type User struct {
		ID   string `validate:"required"`
		Name string `validate:"required"`
	}

	v.Validate(ctx, User{})
	v.Validate(ctx, User{ID: "123"})
	v.Validate(ctx, User{ID: "123", Name: "John"})
}
