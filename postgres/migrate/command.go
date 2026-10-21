package migrate

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Registers the "pgx" database/sql driver the command connects with.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/lcnascimento/go-kit/o11y"
)

const timeout = 5 * time.Minute

var directions = map[string]func(context.Context, *sql.DB, fs.FS) error{
	"up":   Up,
	"down": Down,
}

// Run is the body of a migrate binary: args are its command-line arguments,
// and it exits non-zero when the migrations fail.
func Run(args []string, migrations fs.FS, databaseURL func() (string, error)) {
	if code := run(args, migrations, databaseURL); code != 0 {
		os.Exit(code)
	}
}

func run(args []string, migrations fs.FS, databaseURL func() (string, error)) int {
	ctx := context.Background()

	o11y.MustStart(ctx)

	defer func() {
		if err := o11y.Shutdown(ctx); err != nil {
			onO11yShutdownFailed(ctx, err)
		}
	}()

	signalCtx, stopSignals := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	commandCtx, cancel := context.WithTimeout(signalCtx, timeout)
	defer cancel()

	if err := command(commandCtx, args, migrations, databaseURL); err != nil {
		onCommandFailed(ctx, err)

		return 1
	}

	return 0
}

func command(ctx context.Context, args []string, migrations fs.FS, databaseURL func() (string, error)) error {
	direction := "up"
	if len(args) > 0 {
		direction = args[0]
	}

	migrate, ok := directions[direction]
	if len(args) > 1 || !ok {
		return ErrUsage
	}

	url, err := databaseURL()
	if err != nil {
		return err
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return ErrConnectDatabase.WithCause(err)
	}

	defer func() { _ = db.Close() }()

	return migrate(ctx, db, migrations)
}
