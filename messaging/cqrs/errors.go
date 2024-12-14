package cqrs

import "github.com/lcnascimento/go-kit/errors"

var (
	// ErrBuildRouter is returned when failed to build router.
	ErrBuildRouter errors.CustomError = errors.New("failed to build router").
			WithCode("BUILD_ROUTER_ERROR").
			WithKind(errors.KindInternal)

		// ErrBuildCommandBus is returned when failed to build command bus.
	ErrBuildCommandBus errors.CustomError = errors.New("failed to build command bus").
				WithCode("BUILD_COMMAND_BUS_ERROR").
				WithKind(errors.KindInternal)

	// ErrBuildCommandProcessor is returned when failed to build command processor.
	ErrBuildCommandProcessor errors.CustomError = errors.New("failed to build command processor").
					WithCode("BUILD_COMMAND_PROCESSOR_ERROR").
					WithKind(errors.KindInternal)

	// ErrBuildEventBus is returned when failed to build event bus.
	ErrBuildEventBus errors.CustomError = errors.New("failed to build event bus").
				WithCode("BUILD_EVENT_BUS_ERROR").
				WithKind(errors.KindInternal)

	// ErrBuildEventProcessor is returned when failed to build event processor.
	ErrBuildEventProcessor errors.CustomError = errors.New("failed to build event processor").
				WithCode("BUILD_EVENT_PROCESSOR_ERROR").
				WithKind(errors.KindInternal)

	// ErrAddCommandHandlers is returned when failed to add command handlers into command processor.
	ErrAddCommandHandlers errors.CustomError = errors.New("failed to add command handlers into command processor").
				WithCode("ADD_COMMAND_HANDLERS_ERROR").
				WithKind(errors.KindInternal)

	// ErrAddEventHandlers is returned when failed to add event handlers into event processor.
	ErrAddEventHandlers errors.CustomError = errors.New("failed to add event handlers into event processor").
				WithCode("ADD_EVENT_HANDLERS_ERROR").
				WithKind(errors.KindInternal)

	// ErrRunRouter is returned when failed to run router.
	ErrRunRouter errors.CustomError = errors.New("something went wrong with Message Router").
			WithCode("RUN_ROUTER_ERROR").
			WithKind(errors.KindInternal)

	// ErrSendCommand is returned when failed to send command to the command bus.
	ErrSendCommand errors.CustomError = errors.New("failed to send command to the command bus").
			WithCode("SEND_COMMAND_ERROR").
			WithKind(errors.KindInternal)

	// ErrSendEvent is returned when failed to send event to the event bus.
	ErrSendEvent errors.CustomError = errors.New("failed to send event to the event bus").
			WithCode("SEND_EVENT_ERROR").
			WithKind(errors.KindInternal)
)
