package migrate_test

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/postgres/migrate"
	"github.com/lcnascimento/go-kit/postgres/testdb"
)

var migrations = os.DirFS("testdata")

func TestComponentMigrate_UpAppliesEveryMigration(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.NoError(t, migrate.Up(t.Context(), db, migrations))

	require.True(t, tableExists(t, db, "widgets"))
	require.True(t, tableExists(t, db, "gadgets"))
}

func TestComponentMigrate_RerunIsNoop(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.NoError(t, migrate.Up(t.Context(), db, migrations))
	require.NoError(t, migrate.Up(t.Context(), db, migrations))
}

func TestComponentMigrate_ConcurrentRunsAllSucceed(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	const runs = 4

	errs := make([]error, runs)

	var wg sync.WaitGroup
	for i := range runs {
		wg.Go(func() { errs[i] = migrate.Up(t.Context(), db, migrations) })
	}

	wg.Wait()

	for _, err := range errs {
		require.NoError(t, err)
	}
}

func TestComponentMigrate_AppliedVersionWithoutFileAsksForReset(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.NoError(t, migrate.Up(t.Context(), db, migrations))
	insertUnknownVersion(t, db)

	err := migrate.Up(t.Context(), db, migrations)

	require.ErrorIs(t, err, migrate.ErrUnknownAppliedVersion)
	require.Contains(t, err.Error(), "reset the database")
}

func TestComponentMigrate_DownRevertsTheLatestMigration(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.NoError(t, migrate.Up(t.Context(), db, migrations))
	require.NoError(t, migrate.Down(t.Context(), db, migrations))

	require.False(t, tableExists(t, db, "gadgets"))
	require.True(t, tableExists(t, db, "widgets"))
}

func TestComponentMigrate_DownWithNothingAppliedFails(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.ErrorIs(t, migrate.Down(t.Context(), db, migrations), migrate.ErrNothingToRollBack)
}

func TestComponentMigrate_DownWithAppliedVersionWithoutFileAsksForReset(t *testing.T) {
	t.Parallel()

	db := newDatabase(t)

	require.NoError(t, migrate.Up(t.Context(), db, migrations))
	insertUnknownVersion(t, db)

	err := migrate.Down(t.Context(), db, migrations)

	require.ErrorIs(t, err, migrate.ErrUnknownAppliedVersion)
	require.True(t, tableExists(t, db, "gadgets"))
}

func insertUnknownVersion(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.ExecContext(t.Context(), `INSERT INTO goose_db_version (version_id, is_applied) VALUES (99999, true)`)
	require.NoError(t, err)
}

func tableExists(t *testing.T, db *sql.DB, table string) bool {
	t.Helper()

	var exists bool

	err := db.QueryRowContext(t.Context(), `SELECT to_regclass('public.' || $1) IS NOT NULL`, table).Scan(&exists)
	require.NoError(t, err)

	return exists
}

func newDatabase(t *testing.T) *sql.DB {
	t.Helper()

	if testing.Short() {
		t.Skip("component test")
	}

	return testdb.NewEmpty(t)
}
