package cqrs

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

func (b *broker) onStart(ctx context.Context) {
	logger.Info(ctx, "starting cqrs broker")
}

func (b *broker) onStop(ctx context.Context) {
	logger.Info(ctx, "stopping cqrs broker")
}

func (b *broker) onRunning(ctx context.Context) {
	logger.Info(ctx, "cqrs broker is running")
}

func (b *broker) onAddCommandHandlersError(ctx context.Context, err error) error {
	err = ErrAddCommandHandlers.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (b *broker) onAddEventHandlersError(ctx context.Context, err error) error {
	err = ErrAddEventHandlers.WithCause(err)

	logger.Error(ctx, err)

	return err
}

func (b *broker) onSendCommandError(ctx context.Context, cmd any, err error) error {
	err = ErrSendCommand.WithCause(err)

	logger.Error(ctx, err, log.Any(logKeyCommand, cmd))

	return err
}

func (b *broker) onSendEventError(ctx context.Context, event any, err error) error {
	err = ErrSendEvent.WithCause(err)

	logger.Error(ctx, err, log.Any(logKeyEvent, event))

	return err
}
