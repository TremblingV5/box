package dbtx

type result[R any, E any] struct {
	r R
	e E
}

func newResult[R any, E any](r R, e E) *result[R, E] {
	return &result[R, E]{
		r: r,
		e: e,
	}
}

func (r *result[R, E]) all() (R, E) {
	return r.r, r.e
}

func (r *result[R, E]) left() R {
	return r.r
}

func (r *result[R, E]) right() E {
	return r.e
}
