package cqrs

import (
	"go.opentelemetry.io/otel"

	"github.com/lcnascimento/go-kit/o11y/log"
)

const (
	logKeyCommand            = "command"
	logKeyCommandHandlerName = "command_handler_name"
	logKeyEvent              = "event"
	logKeyEventHandlerName   = "event_handler_name"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/messaging/cqrs"
	logger = log.NewLogger(pkg)
	tracer = otel.Tracer(pkg)
)
