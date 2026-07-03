package util

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lcnascimento/go-kit/o11y/baggage"
)

// ToPointer returns a pointer reference to the given object.
func ToPointer[T any](v T) *T {
	return &v
}

// SafeValue returns the value associated to a pointer.
// If nil, returns the zero value of the given type.
func SafeValue[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

// FakeNow returns a time.Time that is always the same.
// This is useful for testing purposes.
func FakeNow() time.Time {
	return time.Date(2024, 2, 5, 14, 35, 0, 0, time.UTC)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func SortedKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
	})

	return keys
}

func CorrelationID(ctx context.Context) string {
	cID, ok := ctx.Value(baggage.MemberKeyCorrelationID).(string)
	if ok {
		return cID
	}

	return uuid.New().String()
}
