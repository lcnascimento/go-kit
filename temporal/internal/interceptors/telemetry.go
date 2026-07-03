package interceptors

import (
	"context"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger = log.MustNewLogger("github.com/lcnascimento/go-kit/temporal/internal/interceptors")

func onActivityError(ctx context.Context, err error) {
	logger.ErrorBySeverity(ctx, err)
}
