package testdb_test

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/postgres/testdb"
)

var migrations = os.DirFS("../migrate/testdata")

func TestComponentTestdb_NewAppliesTheMigrations(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	var tables int

	err := pool.QueryRow(t.Context(),
		`SELECT count(*) FROM pg_tables WHERE schemaname = 'public' AND tablename IN ('widgets', 'gadgets')`,
	).Scan(&tables)
	require.NoError(t, err)
	require.Equal(t, 2, tables)
}

func TestComponentTestdb_NewIsolatesEachTest(t *testing.T) {
	t.Parallel()

	first := newPool(t)
	second := newPool(t)

	_, err := first.Exec(t.Context(), `INSERT INTO widgets (id) VALUES ('w-1')`)
	require.NoError(t, err)

	var count int

	require.NoError(t, second.QueryRow(t.Context(), `SELECT count(*) FROM widgets`).Scan(&count))
	require.Equal(t, 0, count)
}

// Superuser bypasses row-level security even under FORCE, so a test running
// as one would pass without exercising the policies.
func TestComponentTestdb_NewConnectsWithoutSuperuser(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	var superuser bool

	require.NoError(t, pool.QueryRow(t.Context(), `SELECT rolsuper FROM pg_roles WHERE rolname = current_user`).Scan(&superuser))
	require.False(t, superuser)
}

func TestComponentTestdb_NewEmptyHasNoTables(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("component test")
	}

	db := testdb.NewEmpty(t)

	var tables int

	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM pg_tables WHERE schemaname = 'public'`).Scan(&tables))
	require.Equal(t, 0, tables)
}

func TestComponentTestdb_NewEmptyConnectsWithoutSuperuser(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("component test")
	}

	db := testdb.NewEmpty(t)

	var superuser bool

	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT rolsuper FROM pg_roles WHERE rolname = current_user`).Scan(&superuser))
	require.False(t, superuser)
}

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("component test")
	}

	return testdb.New(t, migrations)
}
