package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/o11y/log"

	"github.com/lcnascimento/go-kit/http/httpserver/util"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if cause := recover(); cause != nil {
				logger.Critical(
					ctx,
					errors.New("%s", fmt.Sprint(cause)).WithCode("PANIC").WithKind(errors.KindCritical),
					log.String("exception.stack", string(debug.Stack())),
				)

				util.WriteError(r.Context(), w, errors.New("unexpected server error")) //nolint:contextcheck // OK
			}
		}(r.Context())

		next.ServeHTTP(w, r)
	})
}
