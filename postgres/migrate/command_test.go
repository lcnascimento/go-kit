package migrate_test

import (
	"database/sql"
	"testing"

	"github.com/peterldowns/pgtestdb"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/errors"

	"github.com/lcnascimento/go-kit/postgres/internal/admin"
	"github.com/lcnascimento/go-kit/postgres/migrate"
)

func TestCommand_UnknownDirectionFailsWithUsage(t *testing.T) {
	t.Parallel()

	err := migrate.Command(t.Context(), []string{"sideways"}, migrations, unreachableURL(t))

	require.ErrorIs(t, err, migrate.ErrUsage)
	require.Equal(t, "down, up", errors.Attributes(err)["directions"])
}

func TestCommand_ExtraArgumentFailsWithUsage(t *testing.T) {
	t.Parallel()

	err := migrate.Command(t.Context(), []string{"up", "now"}, migrations, unreachableURL(t))

	require.ErrorIs(t, err, migrate.ErrUsage)
}

func TestCommand_MissingDatabaseURLFails(t *testing.T) {
	t.Parallel()

	err := migrate.Command(t.Context(), nil, migrations, func() (string, error) { return "", errors.ErrMock })

	require.ErrorIs(t, err, migrate.ErrDatabaseURL)
	require.ErrorIs(t, err, errors.ErrMock)
}

func TestComponentCommand_DefaultsToUp(t *testing.T) {
	t.Parallel()

	url, db := newURL(t)

	require.NoError(t, migrate.Command(t.Context(), nil, migrations, func() (string, error) { return url, nil }))

	require.True(t, tableExists(t, db, "gadgets"))
}

func TestComponentCommand_DownRevertsTheLatestMigration(t *testing.T) {
	t.Parallel()

	url, db := newURL(t)
	databaseURL := func() (string, error) { return url, nil }

	require.NoError(t, migrate.Command(t.Context(), []string{"up"}, migrations, databaseURL))
	require.NoError(t, migrate.Command(t.Context(), []string{"down"}, migrations, databaseURL))

	require.False(t, tableExists(t, db, "gadgets"))
	require.True(t, tableExists(t, db, "widgets"))
}

func unreachableURL(t *testing.T) func() (string, error) {
	t.Helper()

	return func() (string, error) {
		t.Fatal("the database URL must not be read when the arguments are invalid")

		return "", nil
	}
}

func newURL(t *testing.T) (string, *sql.DB) {
	t.Helper()

	if testing.Short() {
		t.Skip("component test")
	}

	conf := admin.Config(t)
	instance := pgtestdb.Custom(t, conf, pgtestdb.NoopMigrator{})

	db, err := instance.Connect()
	require.NoError(t, err)

	t.Cleanup(func() { _ = db.Close() })

	return instance.URL(), db
}
