package command

import (
	"context"

	"github.com/lcnascimento/go-kit/messaging"
	"github.com/lcnascimento/go-kit/messaging/example/event"
	"github.com/lcnascimento/go-kit/o11y/log"
)

// CommandLogMessage is the command that is sent to log a message.
type CommandLogMessage struct {
	Message string
}

// CommandLogMessageHandler is the handler for the CommandLogMessage command.
type CommandLogMessageHandler struct {
	Broker messaging.BrokerCQRS
	Logger *log.Logger
}

// HandlerName returns the name of the handler.
func (h *CommandLogMessageHandler) HandlerName() string {
	return "command.log_message"
}

// NewCommand returns a new instance of the command.
func (h *CommandLogMessageHandler) NewCommand() any {
	return &CommandLogMessage{}
}

// Handle handles the command.
func (h *CommandLogMessageHandler) Handle(ctx context.Context, c any) error {
	command, _ := c.(*CommandLogMessage)
	h.Logger.Info(ctx, command.Message)

	return h.Broker.SendEvent(ctx, &event.EventMessageLogged{})
}
