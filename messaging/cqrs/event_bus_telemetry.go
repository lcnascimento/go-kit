package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/o11y/log"
)

func (b *eventBusBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildEventBus.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *eventBusBuilder) onPublishStart(params cqrs.OnEventSendParams) (context.Context, trace.Span) {
	ctx := params.Message.Context()

	logger.Debug(
		ctx,
		"sending event",
		log.Any(logKeyEvent, params.Event),
	)

	return tracer.Start(ctx, params.EventName, trace.WithSpanKind(trace.SpanKindProducer))
}

func (b *eventBusBuilder) onPublishEnd(ctx context.Context, span trace.Span) {
	logger.Debug(ctx, "event sent")

	span.End()
}
