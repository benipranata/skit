package skit

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// Deref returns the value pointed to by p, or the zero value of T if p is nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr returns the value pointed to by p, or def if p is nil.
func DerefOr[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}
