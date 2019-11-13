package middleware

// Context for Middleware
type Context map[string]interface{}

// Next func
type Next func()

// Middleware func
type Middleware func(Context, Next)

// Merge Middlewares
func Merge(middlewares ...Middleware) Middleware {
	return func(c Context, next Next) {
		for i := len(middlewares) - 1; i > 0; i-- {
			m := middlewares[i]
			n := next
			next = func() {
				m(c, n)
			}
		}
		middlewares[0](c, next)
	}
}
