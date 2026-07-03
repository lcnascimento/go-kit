package propagators

import (
	"context"

	"go.temporal.io/sdk/workflow"

	otel "go.opentelemetry.io/otel/baggage"

	temporal "github.com/lcnascimento/go-kit/temporal/baggage"
)

const headerKey string = "temporal-baggage-header"

type (
	baggagePropagator struct{}
)

func NewBaggagePropagator() workflow.ContextPropagator {
	return &baggagePropagator{}
}

func (s *baggagePropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	bag := otel.FromContext(ctx)

	payload, err := ToPayload(bag)
	if err != nil {
		return err
	}

	writer.Set(headerKey, payload)

	return nil
}

func (s *baggagePropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	bag := temporal.FromContext(ctx)

	payload, err := ToPayload(bag)
	if err != nil {
		return err
	}

	writer.Set(headerKey, payload)

	return nil
}

func (s *baggagePropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	value, ok := reader.Get(headerKey)
	if !ok {
		return ctx, nil
	}

	var bag otel.Baggage
	if err := FromPayload(value, &bag); err == nil {
		ctx = otel.ContextWithBaggage(ctx, bag)
	}

	return ctx, nil
}

func (s *baggagePropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	value, ok := reader.Get(headerKey)
	if !ok {
		return ctx, nil
	}

	var bag temporal.Baggage
	if err := FromPayload(value, &bag); err == nil {
		ctx = temporal.ContextWithBaggage(ctx, bag)
	}

	return ctx, nil
}
