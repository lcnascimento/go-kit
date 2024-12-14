package event

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

// EventMessageLogged is the event that is sent when a message is logged.
type EventMessageLogged struct{}

// EventMessageLoggedHandler is the handler for the EventMessageLogged event.
type EventMessageLoggedHandler struct {
	Logger *log.Logger
}

// HandlerName returns the name of the handler.
func (h *EventMessageLoggedHandler) HandlerName() string {
	return "event.message_logged"
}

// NewEvent returns a new instance of the event.
func (h *EventMessageLoggedHandler) NewEvent() any {
	return &EventMessageLogged{}
}

// Handle handles the event.
func (h *EventMessageLoggedHandler) Handle(ctx context.Context, _ any) error {
	h.Logger.Info(ctx, "message logged")
	return nil
}
