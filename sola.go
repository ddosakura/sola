package sola

import (
	"net/http"

	"github.com/jinzhu/gorm"
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
		// core
		middlewares []Middleware

		// handler
		handlers     map[int]Handler      // not 500
		ErrorHandler func(error, Context) // 500

		// config
		devMode bool

		// orm
		orm map[string]*gorm.DB
	}
)

// New Sola App
func New() *Sola {
	return &Sola{
		middlewares: []Middleware{},
		handlers:    map[int]Handler{},
		orm:         map[string]*gorm.DB{},
	}
}

// Use Middleware
func (s *Sola) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// ServeHTTP to impl http handler
func (s *Sola) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Context{}
	c[CtxSola] = s
	c[CtxRequest] = r
	c[CtxResponse] = w

	h := Merge(s.middlewares...).Handler()
	if err := h(c); err != nil {
		if s.ErrorHandler != nil {
			s.ErrorHandler(err, c)
		} else if s.devMode {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
