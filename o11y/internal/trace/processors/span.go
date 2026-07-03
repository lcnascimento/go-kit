package processors

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
)

type baggageSpanProcessor struct {
	baggageKeysToAddToTags map[string]bool
}

type config struct {
	baggageKeysToAddToTags map[string]bool
}

type Option func(*config)

func WithBaggageKeysToAddToTags(keys ...string) Option {
	return func(cfg *config) {
		for _, key := range keys {
			cfg.baggageKeysToAddToTags[key] = true
		}
	}
}

func BaggageSpanProcessor(opts ...Option) trace.SpanProcessor {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &baggageSpanProcessor{
		baggageKeysToAddToTags: cfg.baggageKeysToAddToTags,
	}
}

func (p *baggageSpanProcessor) OnStart(parentCtx context.Context, s trace.ReadWriteSpan) {
	bag := baggage.FromContext(parentCtx)

	if bag.Len() > 0 {
		attrs := make([]attribute.KeyValue, 0, len(p.baggageKeysToAddToTags))
		for _, member := range bag.Members() {
			if _, ok := p.baggageKeysToAddToTags[member.Key()]; ok {
				attrs = append(attrs, attribute.String(member.Key(), member.Value()))
			}
		}

		s.SetAttributes(attrs...)
	}
}

func (p *baggageSpanProcessor) OnEnd(_ trace.ReadOnlySpan) {}

func (p *baggageSpanProcessor) Shutdown(_ context.Context) error { return nil }

func (p *baggageSpanProcessor) ForceFlush(_ context.Context) error { return nil }
