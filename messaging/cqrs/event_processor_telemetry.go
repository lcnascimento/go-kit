package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/lcnascimento/go-kit/o11y/log"
	"go.opentelemetry.io/otel/trace"
)

func (b *eventProcessorBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildEventProcessor.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *eventProcessorBuilder) onStart(params cqrs.EventProcessorOnHandleParams) (context.Context, trace.Span) {
	ctx := params.Message.Context()

	logger.Debug(
		ctx,
		"event handling started",
		log.String(logKeyEventHandlerName, params.Handler.HandlerName()),
		log.Any(logKeyEvent, params.Event),
	)

	return tracer.Start(ctx, params.Handler.HandlerName(), trace.WithSpanKind(trace.SpanKindConsumer))
}

func (b *eventProcessorBuilder) onEnd(ctx context.Context, span trace.Span, params cqrs.EventProcessorOnHandleParams) {
	logger.Debug(
		ctx,
		"event handling ended",
		log.String(logKeyEventHandlerName, params.Handler.HandlerName()),
	)

	span.End()
}
