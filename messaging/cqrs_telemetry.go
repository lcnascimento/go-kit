package messaging

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger = log.NewLogger("github.com/lcnascimento/go-kit/messaging")

func (b *brokerCQRS) onBuildRouterError(ctx context.Context, err error) error {
	err = ErrBuildRouter.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *brokerCQRS) onBuildCommandBusError(ctx context.Context, err error) error {
	err = ErrBuildCommandBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *brokerCQRS) onBuildCommandProcessorError(ctx context.Context, err error) error {
	err = ErrBuildCommandProcessor.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *brokerCQRS) onBuildEventBusError(ctx context.Context, err error) error {
	err = ErrBuildEventBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *brokerCQRS) onBuildEventProcessorError(ctx context.Context, err error) error {
	err = ErrBuildEventProcessor.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *brokerCQRS) onAddCommandHandlersError(ctx context.Context, err error) error {
	err = ErrAddCommandHandlers.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (b *brokerCQRS) onAddEventHandlersError(ctx context.Context, err error) error {
	err = ErrAddEventHandlers.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (b *brokerCQRS) onStart(ctx context.Context) {
	logger.Info(ctx, "starting cqrs broker")
}

func (b *brokerCQRS) onStop(ctx context.Context) {
	logger.Info(ctx, "stopping cqrs broker")
}

func (b *brokerCQRS) onRunning(ctx context.Context) {
	logger.Info(ctx, "cqrs broker is running")
}

func (b *brokerCQRS) onSendCommandError(ctx context.Context, cmd any, err error) error {
	err = ErrSendCommand.WithCause(err)

	logger.Error(ctx, err, log.Any(logKeyCommand, cmd))

	return err
}

func (b *brokerCQRS) onSendEventError(ctx context.Context, event any, err error) error {
	err = ErrSendEvent.WithCause(err)

	logger.Error(ctx, err, log.Any(logKeyEvent, event))

	return err
}
