package propagation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lcnascimento/go-kit/propagation"
)

func TestContextKeySet(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		t.Run("no value", func(t *testing.T) {
			set := propagation.ContextKeySet{}
			set = set.Add(propagation.ContextKey("foo"))

			assert.Len(t, set, 1)
			assert.Equal(t, true, set[propagation.ContextKey("foo")])
		})

		t.Run("with value", func(t *testing.T) {
			set := propagation.ContextKeySet{}
			set = set.Add(propagation.ContextKey("foo"), "custom-value")

			assert.Len(t, set, 1)
			assert.Equal(t, "custom-value", set[propagation.ContextKey("foo")])
		})
	})

	t.Run("Merge", func(t *testing.T) {
		set1 := propagation.ContextKeySet{
			propagation.ContextKey("foo"): true,
			propagation.ContextKey("bar"): true,
		}

		set2 := propagation.ContextKeySet{
			propagation.ContextKey("foo"): false,
			propagation.ContextKey("baz"): false,
		}

		final := set1.Merge(set2)

		assert.Equal(t, 3, len(final))
		assert.Equal(t, false, final[propagation.ContextKey("foo")])
		assert.Equal(t, true, final[propagation.ContextKey("bar")])
		assert.Equal(t, false, final[propagation.ContextKey("baz")])
	})
}
