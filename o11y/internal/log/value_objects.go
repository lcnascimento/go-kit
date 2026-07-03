//nolint:revive // OK
package log

import (
	"context"
	"log/slog"

	"github.com/lcnascimento/go-kit/o11y/baggage"
)

func baggageAttr(ctx context.Context) slog.Attr {
	bag := baggage.FromContext(ctx)

	attrs := make([]any, 0, len(bag.Members()))

	for _, member := range bag.Members() {
		attrs = append(attrs, slog.String(member.Key(), member.Value()))
	}

	if len(attrs) == 0 {
		return slog.Attr{}
	}

	return slog.Group("bag", attrs...)
}
