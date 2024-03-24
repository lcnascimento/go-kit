package propagation

// ContextKey is a type that should be used to store data into context.
// It is essential to guarantee data context propagation features on logging and distributed tracing.
type ContextKey string

// ContextKeySet is a set of ContextKeys.
// Developers are capable to use the value of this map freely.
// Example: Map a ContextKey to an HTTP Request Param.
type ContextKeySet map[ContextKey]any

// Add creates a new ContextKeySet with the given ContextKey attached.
func (s ContextKeySet) Add(key ContextKey, value ...any) ContextKeySet {
	if len(value) > 0 {
		s[key] = value[0]
	} else {
		s[key] = true
	}

	return s
}

// Merge merges the given ContextKeySet into the existing one and returns a copy of the result.
func (s ContextKeySet) Merge(set ContextKeySet) ContextKeySet {
	for attr := range set {
		s[attr] = set[attr]
	}

	return s
}
