package middlewares

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/lcnascimento/go-kit/o11y/baggage"
)

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cID := r.Header.Get("X-Correlation-Key")
		if cID == "" {
			cID = uuid.New().String()
		}

		r = r.WithContext(baggage.ContextWithCorrelationID(r.Context(), cID))
		w.Header().Set("X-Correlation-Key", cID)

		next.ServeHTTP(w, r)
	})
}
