package sola

import (
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"
)

type (
	// Context for Middleware
	Context interface {
		// Set/Get
		Store() map[string]interface{}
		Origin() Context
		Shadow() Context
		Set(key string, value interface{})
		Get(key string) interface{}

		// API
		Sola() *Sola
		SetCookie(cookie *http.Cookie)
		Request() *http.Request
		Response() http.ResponseWriter

		// Writer
		Blob(code int, contentType string, bs []byte) (err error)
		HTML(code int, data string) error
		String(code int, data string) error
		JSON(code int, data interface{}) error
		File(f File) (err error)

		// Reader
		GetJSON(data interface{}) error

		// Handler
		Handle(code int) Handler
	}
	context struct {
		origin Context
		lock   sync.RWMutex
		store  map[string]interface{}
	}

	// Handler func
	Handler func(Context) error

	// Middleware func
	Middleware func(Handler) Handler

	// C alias of Context
	C Context

	// H alias of Handler
	H Handler

	// M func
	M func(c C, next H) error

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
		handlers: map[int]Handler{
			HandleCodePass: HandlePass,
		},
		orm: map[string]*gorm.DB{},
	}
}

// Dev Mode
func (s *Sola) Dev() {
	s.devMode = true
}

// Use Middleware
func (s *Sola) Use(ms ...Middleware) {
	s.middlewares = append(s.middlewares, ms...)
}

// ServeHTTP to impl http handler
func (s *Sola) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext()
	c.Set(CtxSola, s)
	c.Set(CtxRequest, r)
	c.Set(CtxResponse, w)

	h := Merge(s.middlewares...).Handler()
	if err := h(c); err != nil {
		if s.ErrorHandler != nil {
			s.ErrorHandler(err, c)
		} else if s.devMode {
			if ie, ok := err.(*InternalError); ok {
				ie.Write(w)
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
