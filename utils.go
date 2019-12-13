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

// === N/M 互转 & Must ===

// M to Middleware
func (m M) M() Middleware {
	return func(next Handler) Handler {
		return func(c Context) error {
			return m(C(c), H(next))
		}
	}
}

// Must not nil
func (m M) Must(not Handler) Middleware {
	return m.M().Must(not)
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

// Must not nil
func (h Handler) Must(not Handler) Handler {
	if h == nil {
		return not
	}
	return h
}

// H wrap Middleware to Handler
func (m Middleware) H() Handler {
	return m(errM2H.Handler())
}

// Must not nil
func (m Middleware) Must(not Handler) Middleware {
	return func(next Handler) Handler {
		next = next.Must(not)
		if m == nil {
			return next
		}
		return m(next)
	}
}

// === Adapter for net/http ===

// FromHandlerFunc func(http.ResponseWriter, *http.Request)
func FromHandlerFunc(h func(http.ResponseWriter, *http.Request)) Handler {
	return func(c Context) error {
		h(c.Response(), c.Request())
		return nil
	}
}

// FromHandler http.Handler
func FromHandler(h http.Handler) Handler {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// ToHandlerFunc func(http.ResponseWriter, *http.Request)
func (h Handler) ToHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := newContext()
		c.Set(CtxRequest, r)
		c.Set(CtxResponse, w)
	}
}

// ToHandlerFunc func(http.ResponseWriter, *http.Request)
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ToHandlerFunc()(w, r)
}

// ToHandler http.Handler
func (h Handler) ToHandler() http.Handler {
	return h
}
