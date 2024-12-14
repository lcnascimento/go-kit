package cqrs

import "context"

func (b *eventBusBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildEventBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}
