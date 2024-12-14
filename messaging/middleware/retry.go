package middleware

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cenkalti/backoff/v3"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger = log.NewLogger("github.com/lcnascimento/go-kit/messaging/middleware/retry")

const (
	multiplier      = float64(2)
	initialInterval = 1 * time.Second
	maxInterval     = 2 * time.Minute
	maxRetries      = 100
)

var errMaxRetriesReached = errors.New("discarding message after reaching max retries attempts").
	WithKind(errors.KindInternal).
	WithCode("MAX_RETRIES_REACHED_ERROR")

// Retry is a middleware that retries a message if it fails.
func Retry() message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			producedMessages, err := h(msg)
			if err == nil {
				return producedMessages, nil
			}

			expBackoff := backoff.NewExponentialBackOff()
			expBackoff.InitialInterval = initialInterval
			expBackoff.MaxInterval = maxInterval
			expBackoff.Multiplier = multiplier

			ctx := msg.Context()

			retryNum := 1
			expBackoff.Reset()

			for {
				waitTime := expBackoff.NextBackOff()
				select {
				case <-ctx.Done():
					return producedMessages, err
				case <-time.After(waitTime):
					// go on
				}

				logger.Info(
					ctx, "retrying message",
					log.String("message_id", msg.UUID),
					log.Int("retry_num", retryNum),
					log.String("delay", waitTime.String()),
				)

				producedMessages, err = h(msg)
				if err == nil {
					return producedMessages, nil
				}

				retryNum++
				if retryNum > maxRetries {
					break
				}
			}

			logger.Critical(ctx, errMaxRetriesReached, log.String("message_id", msg.UUID))

			return producedMessages, nil
		}
	}
}
