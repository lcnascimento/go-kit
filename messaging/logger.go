package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/o11y/log"
)

// BaggageFieldPrefix is the prefix for baggage fields.
const BaggageFieldPrefix = "baggage."

// WatermillLogger is a logger adapter for the Watermill library.
type WatermillLogger struct {
	fields watermill.LogFields
}

// NewWatermillLogger creates a new WatermillLogger instance.
func NewWatermillLogger() watermill.LoggerAdapter {
	return &WatermillLogger{fields: make(watermill.LogFields)}
}

// Error logs an error message with the given fields.
func (l *WatermillLogger) Error(_ string, err error, fields watermill.LogFields) {
	ctx, attrs := l.buildContextAndAttrs(fields)
	logger.Error(ctx, err, attrs...)
}

// Info logs an info message with the given fields.
func (l *WatermillLogger) Info(msg string, fields watermill.LogFields) {
	ctx, attrs := l.buildContextAndAttrs(fields)
	logger.Info(ctx, msg, attrs...)
}

// Debug logs a debug message with the given fields.
func (l *WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	ctx, attrs := l.buildContextAndAttrs(fields)
	logger.Debug(ctx, msg, attrs...)
}

// Trace logs a trace message with the given fields.
func (l *WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	ctx, attrs := l.buildContextAndAttrs(fields)
	logger.Trace(ctx, msg, attrs...)
}

// With adds fields to be logged in all log messages.
func (l *WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	l.fields = l.fields.Add(fields)
	return l
}

func (l *WatermillLogger) buildContextAndAttrs(fields watermill.LogFields) (context.Context, []slog.Attr) {
	ctx := context.Background()
	bag, err := baggage.New()
	if err != nil {
		_ = l.onCreateBaggageError(ctx, err)
		return ctx, []slog.Attr{}
	}

	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		if !strings.HasPrefix(k, BaggageFieldPrefix) {
			attrs = append(attrs, log.Any(k, v))
			continue
		}

		member, err := baggage.NewMember(strings.TrimPrefix(k, BaggageFieldPrefix), fmt.Sprintf("%v", v))
		if err != nil {
			_ = l.onCreateBaggageMemberError(ctx, err)
			continue
		}

		bag, err = bag.SetMember(member)
		if err != nil {
			_ = l.onSetBaggageMemberError(ctx, err)
			continue
		}
	}

	if bag.Len() > 0 {
		ctx = baggage.ContextWithBaggage(ctx, bag)
	}

	return ctx, attrs
}
