package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/util"
	"github.com/lcnascimento/go-kit/validator"
)

func TestCanonicalID(t *testing.T) {
	type payload struct {
		ID string `validate:"canonical_id=note"`
	}

	v, err := validator.New()
	require.NoError(t, err)

	t.Run("accepts a string with the expected prefix", func(t *testing.T) {
		err := v.Struct(payload{ID: "note_01920000-0000-7000-8000-000000000000"})

		assert.NoError(t, err)
	})

	for name, id := range map[string]string{
		"wrong prefix":      "user_01920000-0000-7000-8000-000000000000",
		"missing prefix":    "01920000-0000-7000-8000-000000000000",
		"malformed uuid":    "note_not-a-uuid",
		"empty":             "",
		"nil uuid":          "note_00000000-0000-0000-0000-000000000000",
		"prefix as a match": "notebook_01920000-0000-7000-8000-000000000000",
	} {
		t.Run("rejects a string with "+name, func(t *testing.T) {
			err := v.Struct(payload{ID: id})

			assert.ErrorIs(t, err, errors.ErrInvalidInput)
			assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
		})
	}

	t.Run("accepts a util.CanonicalID with the expected prefix", func(t *testing.T) {
		type typed struct {
			ID util.CanonicalID `validate:"canonical_id=note"`
		}

		err := v.Struct(typed{ID: util.NewCanonicalID("note")})

		assert.NoError(t, err)
	})

	t.Run("accepts a type embedding util.CanonicalID with the expected prefix", func(t *testing.T) {
		err := v.Struct(embedding{ID: noteID{util.NewCanonicalID("note")}})

		assert.NoError(t, err)
	})

	t.Run("rejects a util.CanonicalID with the wrong prefix", func(t *testing.T) {
		type typed struct {
			ID util.CanonicalID `validate:"canonical_id=note"`
		}

		err := v.Struct(typed{ID: util.NewCanonicalID("user")})

		assert.ErrorIs(t, err, errors.ErrInvalidInput)
		assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
	})

	t.Run("rejects a zero util.CanonicalID", func(t *testing.T) {
		type typed struct {
			ID util.CanonicalID `validate:"canonical_id=note"`
		}

		err := v.Struct(typed{})

		assert.ErrorIs(t, err, errors.ErrInvalidInput)
		assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
	})

	t.Run("rejects a type embedding a zero util.CanonicalID", func(t *testing.T) {
		err := v.Struct(embedding{})

		assert.ErrorIs(t, err, errors.ErrInvalidInput)
		assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
	})

	t.Run("rejects a type embedding util.CanonicalID with the wrong prefix", func(t *testing.T) {
		err := v.Struct(embedding{ID: noteID{util.NewCanonicalID("user")}})

		assert.ErrorIs(t, err, errors.ErrInvalidInput)
		assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
	})

	t.Run("rejects a field of an unsupported type", func(t *testing.T) {
		type unsupported struct {
			ID int `validate:"canonical_id=note"`
		}

		err := v.Struct(unsupported{ID: 1})

		assert.ErrorIs(t, err, errors.ErrInvalidInput)
		assert.Equal(t, []string{"[ID] must be a valid canonical id with prefix note"}, errors.Reasons(err))
	})
}

type noteID struct{ util.CanonicalID }

type embedding struct {
	ID noteID `validate:"canonical_id=note"`
}
