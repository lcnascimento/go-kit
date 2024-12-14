package cqrs

import "context"

func (b *commandBusBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildCommandBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}
