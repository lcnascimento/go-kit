package util

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
