package sola

// Merge Middlewares
func Merge(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			n := next
			next = func(c Context) error {
				return m(n)(c)
			}
		}
		return next
	}
}

// Handler of Middleware
func (m Middleware) Handler() Handler {
	return m(nil)
}

// MergeFunc for Middlewares & Handler
func MergeFunc(h Handler, middlewares ...Middleware) Handler {
	return Merge(middlewares...)(h)
}
