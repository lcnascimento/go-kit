package cqrs

import "context"

func (b *routerBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildRouter.WithCause(err)

	logger.Critical(ctx, err)

	return err
}
