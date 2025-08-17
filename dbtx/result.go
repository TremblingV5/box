package dbtx

// result is a generic struct that holds two values of potentially different types
// It's typically used to hold a result value and an error
type result[R any, E any] struct {
	r R
	e E
}

// newResult creates a new result instance with the provided values
func newResult[R any, E any](r R, e E) *result[R, E] {
	return &result[R, E]{
		r: r,
		e: e,
	}
}

// all returns both values stored in the result
func (r *result[R, E]) all() (R, E) {
	return r.r, r.e
}

// left returns the first value (typically the result) stored in the result
func (r *result[R, E]) left() R {
	return r.r
}

// right returns the second value (typically the error) stored in the result
func (r *result[R, E]) right() E {
	return r.e
}
