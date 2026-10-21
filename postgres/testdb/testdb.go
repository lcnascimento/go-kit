package testdb

import (
	"context"
	"database/sql"
	"io/fs"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/common"

	"github.com/lcnascimento/go-kit/postgres/internal/admin"
	"github.com/lcnascimento/go-kit/postgres/migrate"
)

// New connects to a fresh database cloned from a template with migrations
// applied, as a role without superuser so row-level security applies.
func New(t *testing.T, migrations fs.FS) *pgxpool.Pool {
	t.Helper()

	instance := pgtestdb.Custom(t, admin.Config(t), migrator{migrations: migrations})

	pool, err := pgxpool.New(t.Context(), instance.URL())
	if err != nil {
		t.Fatalf("testdb: connect to %s: %v", instance.Database, err)
	}

	t.Cleanup(pool.Close)

	return pool
}

func NewEmpty(t *testing.T) *sql.DB {
	t.Helper()

	return pgtestdb.New(t, admin.Config(t), pgtestdb.NoopMigrator{})
}

type migrator struct {
	migrations fs.FS
}

func (m migrator) Hash() (string, error) {
	return common.HashDirs(m.migrations, "*.sql", ".")
}

//nolint:gocritic // hugeParam: the signature is pgtestdb.Migrator's.
func (m migrator) Migrate(ctx context.Context, db *sql.DB, _ pgtestdb.Config) error {
	return migrate.Up(ctx, db, m.migrations)
}
