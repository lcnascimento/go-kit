package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/lcnascimento/go-kit/messaging"
	messagingMid "github.com/lcnascimento/go-kit/messaging/middleware"
)

type routerBuilder struct {
	logger watermill.LoggerAdapter
}

func newRouterBuilder(logger watermill.LoggerAdapter) *routerBuilder {
	return &routerBuilder{
		logger: logger,
	}
}

func (b *routerBuilder) build(ctx context.Context) (*message.Router, error) {
	var err error

	router, err := message.NewRouter(message.RouterConfig{}, b.logger)
	if err != nil {
		return nil, b.onBuildError(ctx, err)
	}

	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(middleware.CorrelationID)
	router.AddMiddleware(messagingMid.WithBaggage(messaging.BaggageFieldPrefix))
	router.AddMiddleware(messagingMid.Retry())

	return router, nil
}
