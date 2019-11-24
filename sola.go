package sola

import (
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"
)

type (
	// Context for Middleware
	Context interface {
		// set/get
		Set(key string, value interface{})
		Get(key string) interface{}

		// api
		Sola() *Sola
		SetCookie(cookie *http.Cookie)
		Request() *http.Request
		Response() http.ResponseWriter

		// Writer
		Blob(code int, contentType string, bs []byte) (err error)
		HTML(code int, data string) error
		String(code int, data string) error
		JSON(code int, data interface{}) error

		// handler
		Handle(code int) Handler
	}
	context struct {
		lock  sync.RWMutex
		store map[string]interface{}
	}

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
	c := &context{
		store: map[string]interface{}{},
	}
	c.Set(CtxSola, s)
	c.Set(CtxRequest, r)
	c.Set(CtxResponse, w)

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
