package baggage

import (
	"go.temporal.io/sdk/workflow"

	otel "go.opentelemetry.io/otel/baggage"

	"github.com/lcnascimento/go-kit/o11y/baggage"
)

type baggageContextKeyType int

const baggageKey baggageContextKeyType = iota

func FromContext(ctx workflow.Context) otel.Baggage {
	empty, _ := New()

	state, ok := ctx.Value(baggageKey).(baggageState)
	if !ok {
		return empty
	}

	if state.bag == nil {
		return empty
	}

	return *state.bag
}

func ContextWithBaggage(parent workflow.Context, b otel.Baggage) workflow.Context {
	return contextWithBaggage(parent, &b)
}

func ContextWithoutBaggage(parent workflow.Context) workflow.Context {
	return contextWithBaggage(parent, nil)
}

func ContextWithCorrelationID(parent workflow.Context, correlationID string) workflow.Context {
	bag := FromContext(parent)

	member, err := otel.NewMember(string(baggage.MemberKeyCorrelationID), correlationID)
	if err != nil {
		return parent
	}

	bag, err = bag.SetMember(member)
	if err != nil {
		return parent
	}

	return contextWithBaggage(parent, &bag)
}

func ContextWithMembers(parent workflow.Context, members ...otel.Member) workflow.Context {
	var err error

	bag := FromContext(parent)

	for _, member := range members {
		bag, err = bag.SetMember(member)
		if err != nil {
			return parent
		}
	}

	return contextWithBaggage(parent, &bag)
}

func contextWithBaggage(parent workflow.Context, b *otel.Baggage) workflow.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.bag = b
	ctx := workflow.WithValue(parent, baggageKey, s)

	return ctx
}
