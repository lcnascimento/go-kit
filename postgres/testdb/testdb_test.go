package testdb_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/postgres/testdb"
)

const roleSkipsRowLevelSecurity = `SELECT rolsuper OR rolbypassrls FROM pg_roles WHERE rolname = current_user`

var migrations = os.DirFS("testdata")

func TestComponentTestdb_NewAppliesTheMigrations(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	var tables int

	err := pool.QueryRow(t.Context(),
		`SELECT count(*) FROM pg_tables WHERE schemaname = 'public' AND tablename IN ('widgets', 'hidden_widgets')`,
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

func TestComponentTestdb_NewEnforcesForcedRowLevelSecurity(t *testing.T) {
	t.Parallel()

	pool := newPool(t)

	var visible int

	require.NoError(t, pool.QueryRow(t.Context(), `SELECT count(*) FROM hidden_widgets`).Scan(&visible))
	require.Equal(t, 0, visible)

	var skips bool

	require.NoError(t, pool.QueryRow(t.Context(), roleSkipsRowLevelSecurity).Scan(&skips))
	require.False(t, skips)
}

func TestComponentTestdb_NewEmptyHasNoTables(t *testing.T) {
	t.Parallel()

	db := newEmpty(t)

	var tables int

	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM pg_tables WHERE schemaname = 'public'`).Scan(&tables))
	require.Equal(t, 0, tables)
}

func TestComponentTestdb_NewEmptyRoleIsSubjectToRowLevelSecurity(t *testing.T) {
	t.Parallel()

	db := newEmpty(t)

	var skips bool

	require.NoError(t, db.QueryRowContext(t.Context(), roleSkipsRowLevelSecurity).Scan(&skips))
	require.False(t, skips)
}

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("component test")
	}

	return testdb.New(t, migrations)
}

func newEmpty(t *testing.T) *sql.DB {
	t.Helper()

	if testing.Short() {
		t.Skip("component test")
	}

	return testdb.NewEmpty(t)
}
