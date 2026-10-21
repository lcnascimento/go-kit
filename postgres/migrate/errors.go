package migrate

import (
	"maps"
	"slices"
	"strings"

	"github.com/lcnascimento/go-kit/errors"
)

var (
	ErrUnknownAppliedVersion = errors.New(
		"the database has migrations applied that this checkout does not have (switched branches?): reset the database",
	).WithCode("UNKNOWN_APPLIED_MIGRATION")

	ErrNothingToRollBack = errors.New("no migration is applied: there is nothing to roll back").
				WithCode("NOTHING_TO_ROLL_BACK")

	ErrMigrate = errors.New("failed to run migrations")

	ErrUsage = errors.New("usage: migrate [direction], with direction one of the listed ones").
			WithCode("MIGRATIONS_USAGE").
			WithAttribute("directions", strings.Join(slices.Sorted(maps.Keys(directions)), ", "))

	ErrConnectDatabase = errors.New("failed to connect to the database: is Postgres up and the database URL right?")
)
