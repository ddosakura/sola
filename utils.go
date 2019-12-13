package sola

import (
	"net/http"
)

// OriginContext Middleware change ctx to it's origin
func OriginContext(next Handler) Handler {
	return func(c Context) error {
		if origin := c.Origin(); origin != nil {
			return next(origin)
		}
		return next(c)
	}
}

// Adapter for net/http
func (h Handler) Adapter() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, e error) {
		c := newContext()
		c.Set(CtxRequest, r)
		c.Set(CtxResponse, w)
	}
}

// M wrap Handler to Middleware
func (h Handler) M() Middleware {
	return func(next Handler) Handler {
		if next == nil {
			return h
		}
		return func(c Context) error {
			if e := h(c); e != nil {
				return e
			}
			return next(c)
		}
	}
}

// H wrap Middleware to Handler
func (m Middleware) H() Handler {
	return m(nil)
}

// PASS Handler
var PASS = Handler(func(c Context) error {
	return nil
})

// Must not nil
func (h Handler) Must() Handler {
	if h == nil {
		return PASS
	}
	return h
}

// MustC not nil
func (h Handler) MustC(c Context) error {
	if h == nil {
		return nil
	}
	return h(c)
}
