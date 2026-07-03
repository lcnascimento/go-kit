package middlewares

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"

	"github.com/lcnascimento/go-kit/http/httpserver/util"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if err := recover(); err != nil {
				logger.CriticalMessage(
					ctx,
					"panic recovered",
					log.Any("exception.message", err),
					log.String("exception.stack", string(debug.Stack())),
				)

				util.WriteError(r.Context(), w, errors.New("unexpected server error")) //nolint:contextcheck // OK
			}
		}(r.Context())

		next.ServeHTTP(w, r)
	})
}
