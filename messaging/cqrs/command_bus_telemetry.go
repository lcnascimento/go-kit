package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/lcnascimento/go-kit/o11y/log"
	"go.opentelemetry.io/otel/trace"
)

func (b *commandBusBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildCommandBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *commandBusBuilder) onSendStart(params cqrs.CommandBusOnSendParams) (context.Context, trace.Span) {
	ctx := params.Message.Context()

	logger.Debug(
		ctx,
		"sending command",
		log.Any(logKeyCommand, params.Command),
	)

	return tracer.Start(ctx, params.CommandName, trace.WithSpanKind(trace.SpanKindProducer))
}

func (b *commandBusBuilder) onSendEnd(ctx context.Context, span trace.Span) {
	logger.Debug(ctx, "command sent")

	span.End()
}
