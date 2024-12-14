package messaging

import "context"

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
