package messaging

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger = log.NewLogger("github.com/lcnascimento/go-kit/messaging")

func (l *WatermillLogger) onCreateBaggageError(ctx context.Context, err error) error {
	err = ErrCreateBaggage.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (l *WatermillLogger) onCreateBaggageMemberError(ctx context.Context, err error) error {
	err = ErrCreateBaggageMember.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (l *WatermillLogger) onSetBaggageMemberError(ctx context.Context, err error) error {
	err = ErrSetBaggageMember.WithCause(err)

	logger.Error(ctx, err)

	return err
}
