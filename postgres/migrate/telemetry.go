package migrate

import (
	"context"

	"github.com/pressly/goose/v3"

	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/postgres/migrate"
	logger = log.MustNewLogger(pkg)
)

func onMigrated(ctx context.Context, result *goose.MigrationResult) {
	if result.Error != nil {
		return
	}

	logger.Info(ctx, "migration "+result.Direction,
		log.String("direction", result.Direction),
		log.Int("version", int(result.Source.Version)),
		log.String("file", result.Source.Path),
		log.String("duration", result.Duration.String()),
	)
}

func onCommandFailed(ctx context.Context, err error) {
	logger.Error(ctx, err)
}

func onO11yShutdownFailed(ctx context.Context, err error) {
	logger.Error(ctx, err)
}
