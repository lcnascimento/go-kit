package validator

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLogger("github.com/lcnascimento/go-kit/validator")
}

func onInvalidStruct(ctx context.Context, s any, err error) error {
	err = ErrInvalidStruct(err)
	logger.Error(ctx, err, log.Any("struct", s))
	return err
}
