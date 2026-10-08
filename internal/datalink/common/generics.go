package common

// Ptr returns a pointer to the given value of any type.
func Ptr[T any](v T) *T {
	return &v
}

// Deref safely dereferences a pointer, returning the fallback value if the pointer is nil.
func Deref[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// DerefOrZero safely dereferences a pointer, returning the type's zero value if the pointer is nil.
func DerefOrZero[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
