package validator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lcnascimento/go-kit/errors"
	"github.com/lcnascimento/go-kit/validator"
)

type noteContent struct {
	Title string `json:"title" validate:"required,max=120"`
}

type embeddedNote struct {
	noteContent
	ID string `json:"id"`
}

type embeddedPointerNote struct {
	*noteContent
	ID string `json:"id"`
}

type nestedNote struct {
	Content noteContent `json:"content"`
	ID      string      `json:"id"`
}

type taggedEmbeddedNote struct {
	noteContent `json:"content"`
	ID          string `json:"id"`
}

type deepNote struct {
	Note embeddedNote `json:"note"`
}

func TestWithJSONFieldNames(t *testing.T) {
	v, err := validator.New(validator.WithJSONFieldNames())
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		payload any
		reason  string
	}{
		"embedded struct contributes no segment": {
			payload: embeddedNote{ID: "1"},
			reason:  "[title] is a required field",
		},
		"embedded struct contributes no segment to a parameterized rule": {
			payload: embeddedNote{noteContent: noteContent{Title: strings.Repeat("a", 121)}, ID: "1"},
			reason:  "[title] must be a maximum of 120 characters in length",
		},
		"embedded pointer contributes no segment": {
			payload: embeddedPointerNote{noteContent: &noteContent{}, ID: "1"},
			reason:  "[title] is a required field",
		},
		"named nested struct keeps its segment": {
			payload: nestedNote{ID: "1"},
			reason:  "[content.title] is a required field",
		},
		"embedded struct with a json name keeps its segment": {
			payload: taggedEmbeddedNote{ID: "1"},
			reason:  "[content.title] is a required field",
		},
		"embedded struct inside a named struct contributes no segment": {
			payload: deepNote{Note: embeddedNote{ID: "1"}},
			reason:  "[note.title] is a required field",
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := v.Struct(tc.payload)

			assert.ErrorIs(t, err, errors.ErrInvalidInput)
			assert.Equal(t, []string{tc.reason}, errors.Reasons(err))
		})
	}
}
