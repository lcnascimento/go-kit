package baggage

import (
	otel "go.opentelemetry.io/otel/baggage"
)

var (
	NewMember = otel.NewMember
	New       = otel.New
)

type (
	Member   = otel.Member
	Baggage  = otel.Baggage
	Property = otel.Property
)

type baggageState struct {
	bag *Baggage
}
