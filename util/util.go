package util

import "time"

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
