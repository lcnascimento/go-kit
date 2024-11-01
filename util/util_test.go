package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lcnascimento/go-kit/util"
)

func TestToPointer(t *testing.T) {
	t.Run("when value is String", func(t *testing.T) {
		v := "test"
		assert.Equal(t, &v, util.ToPointer(v))
	})

	t.Run("when value is Int", func(t *testing.T) {
		v := 1
		assert.Equal(t, &v, util.ToPointer(v))
	})

	t.Run("when value is Bool", func(t *testing.T) {
		v := true
		assert.Equal(t, &v, util.ToPointer(v))
	})
}

func TestSafeValue(t *testing.T) {
	t.Run("when value is nil", func(t *testing.T) {
		var v *string
		assert.Equal(t, "", util.SafeValue(v))
	})

	t.Run("when value is String", func(t *testing.T) {
		v := "test"
		assert.Equal(t, v, util.SafeValue(&v))
	})

	t.Run("when value is Int", func(t *testing.T) {
		v := 1
		assert.Equal(t, v, util.SafeValue(&v))
	})

	t.Run("when value is Bool", func(t *testing.T) {
		v := true
		assert.Equal(t, v, util.SafeValue(&v))
	})
}
