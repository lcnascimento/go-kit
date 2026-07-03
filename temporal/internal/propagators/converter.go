package propagators

import (
	"go.opentelemetry.io/otel/baggage"
	"go.temporal.io/api/common/v1"
)

func ToPayload(bag baggage.Baggage) (*common.Payload, error) {
	data := bag.String()

	return &common.Payload{Data: []byte(data)}, nil
}

func FromPayload(payload *common.Payload, bag *baggage.Baggage) (err error) {
	*bag, err = baggage.Parse(string(payload.GetData()))

	return err
}
