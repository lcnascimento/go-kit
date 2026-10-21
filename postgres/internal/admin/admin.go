package admin

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/peterldowns/pgtestdb"

	// Registers the "pgx" database/sql driver pgtestdb connects with.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/lcnascimento/go-kit/env"
)

const defaultURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" //nolint:gosec // local docker-compose credentials

const pingTimeout = 5 * time.Second

// Config is the reachable admin connection of POSTGRES_ADMIN_URL, failing t
// when Postgres does not answer.
func Config(t *testing.T) pgtestdb.Config {
	t.Helper()

	conf := parse(t)

	if err := ping(t.Context(), &conf); err != nil {
		t.Fatalf("postgres is not reachable at POSTGRES_ADMIN_URL (%s:%s): %v\n"+
			"Start it, point POSTGRES_ADMIN_URL at it, or run the tests with -short to skip the component tests.",
			conf.Host, conf.Port, err)
	}

	return conf
}

func parse(t *testing.T) pgtestdb.Config {
	t.Helper()

	raw := env.Get("POSTGRES_ADMIN_URL", env.WithDefaultValue(defaultURL))

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("testdb: POSTGRES_ADMIN_URL is not a URL: %v", err)
	}

	password, _ := u.User.Password()

	port := u.Port()
	if port == "" {
		port = "5432"
	}

	return pgtestdb.Config{
		DriverName: "pgx",
		Host:       u.Hostname(),
		Port:       port,
		User:       u.User.Username(),
		Password:   password,
		Database:   strings.TrimPrefix(u.Path, "/"),
		Options:    u.RawQuery,
	}
}

func ping(ctx context.Context, conf *pgtestdb.Config) error {
	db, err := conf.Connect()
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	return db.PingContext(ctx)
}
