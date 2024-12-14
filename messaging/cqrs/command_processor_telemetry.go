package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.opentelemetry.io/otel/trace"

	"github.com/lcnascimento/go-kit/o11y/log"
)

func (b *commandProcessorBuilder) onBuildError(ctx context.Context, err error) error {
	err = ErrBuildCommandProcessor.WithCause(err)

	logger.Critical(ctx, err)

	return err
}

func (b *commandProcessorBuilder) onCommandHandlingStart(params cqrs.CommandProcessorOnHandleParams) (context.Context, trace.Span) {
	ctx := params.Message.Context()

	logger.Debug(
		ctx,
		"command handling started",
		log.String(logKeyCommandHandlerName, params.Handler.HandlerName()),
		log.Any(logKeyCommand, params.Command),
	)

	return tracer.Start(ctx, params.Handler.HandlerName(), trace.WithSpanKind(trace.SpanKindConsumer))
}

func (b *commandProcessorBuilder) onCommandHandlingEnded(ctx context.Context, span trace.Span, params cqrs.CommandProcessorOnHandleParams) {
	logger.Debug(
		ctx,
		"command handling ended",
		log.String(logKeyCommandHandlerName, params.Handler.HandlerName()),
	)

	span.End()
}
