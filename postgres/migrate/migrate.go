package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"slices"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"github.com/pressly/goose/v3/lock"

	"github.com/lcnascimento/go-kit/errors"
)

func Up(ctx context.Context, db *sql.DB, migrations fs.FS) error {
	provider, err := newProvider(ctx, db, migrations)
	if err != nil {
		return err
	}

	results, err := provider.Up(ctx)
	for _, result := range results {
		onMigrated(ctx, result)
	}

	if err != nil {
		return ErrMigrate.WithCause(err)
	}

	return nil
}

func Down(ctx context.Context, db *sql.DB, migrations fs.FS) error {
	provider, err := newProvider(ctx, db, migrations)
	if err != nil {
		return err
	}

	result, err := provider.Down(ctx)
	if errors.Is(err, goose.ErrNoNextVersion) {
		return ErrNothingToRollBack
	}

	if result != nil {
		onMigrated(ctx, result)
	}

	if err != nil {
		return ErrMigrate.WithCause(err)
	}

	return nil
}

func newProvider(ctx context.Context, db *sql.DB, migrations fs.FS) (*goose.Provider, error) {
	store, err := database.NewStore(database.DialectPostgres, goose.DefaultTablename)
	if err != nil {
		return nil, ErrMigrate.WithCause(err)
	}

	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, ErrMigrate.WithCause(err)
	}

	provider, err := goose.NewProvider("", db, migrations, goose.WithStore(store), goose.WithSessionLocker(locker))
	if err != nil {
		return nil, ErrMigrate.WithCause(err)
	}

	if err := checkAppliedVersionsAreKnown(ctx, db, store, provider.ListSources()); err != nil {
		return nil, err
	}

	return provider, nil
}

// goose v3 silently skips an applied version with no file, which would leave
// another branch's schema behind unnoticed.
func checkAppliedVersionsAreKnown(ctx context.Context, db *sql.DB, store database.Store, sources []*goose.Source) error {
	extender, ok := store.(database.StoreExtender)
	if !ok {
		return ErrMigrate.WithAttribute("reason", "goose store cannot check whether its table exists")
	}

	exists, err := extender.TableExists(ctx, db)
	if err != nil {
		return ErrMigrate.WithCause(err)
	}

	if !exists {
		return nil
	}

	applied, err := store.ListMigrations(ctx, db)
	if err != nil {
		return ErrMigrate.WithCause(err)
	}

	var unknown []int64

	for _, m := range applied {
		// Version 0 is the row goose inserts when it creates its table.
		known := m.Version == 0 || slices.ContainsFunc(sources, func(s *goose.Source) bool { return s.Version == m.Version })
		if m.IsApplied && !known {
			unknown = append(unknown, m.Version)
		}
	}

	if len(unknown) > 0 {
		return ErrUnknownAppliedVersion.WithAttribute("versions", fmt.Sprint(unknown))
	}

	return nil
}
