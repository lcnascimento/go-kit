package baggage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/o11y/baggage"
)

func TestContextWithAttributes(t *testing.T) {
	ctx := baggage.ContextWithMembers(context.Background(), baggage.NewMember("key", "value"))

	bag := baggage.FromContext(ctx)

	require.Equal(t, bag.Member("key").Value(), "value")
}

func TestContextWithCorrelationID(t *testing.T) {
	cID := uuid.New().String()
	ctx := baggage.ContextWithCorrelationID(context.Background(), cID)

	require.Equal(t, baggage.FromContext(ctx).Member(baggage.MemberKeyCorrelationID).Value(), cID)
}
