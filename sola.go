package sola

import (
	"net/http"
)

type (
	// Context for Middleware
	Context map[string]interface{}

	// Handler func
	Handler func(Context) error

	// Middleware func
	Middleware func(Handler) Handler

	// Sola App
	Sola struct {
		middlewares  []Middleware
		ErrorHandler func(error, Context)
	}
)

// New Sola App
func New() *Sola {
	return &Sola{
		middlewares: []Middleware{},
	}
}

// Use Middleware
func (s *Sola) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// ServeHTTP to impl http handler
func (s *Sola) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Context{}
	c[CtxRequest] = r
	c[CtxResponse] = w

	h := Merge(s.middlewares...).Handler()
	if err := h(c); err != nil {
		if s.ErrorHandler != nil {
			s.ErrorHandler(err, c)
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
