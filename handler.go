package sola

import (
	"net/http"
)

// SetHandler in Sola
func (s *Sola) SetHandler(code int, h Handler) {
	s.handlers[code] = h
}

// Handle Selector
func (c *context) Handle(code int) Handler {
	fn := c.Sola().handlers[code]
	if fn == nil {
		switch code {
		case http.StatusBadRequest:
			fn = func(c Context) error {
				return c.String(code, "Bad Request")
			}
		case http.StatusUnauthorized:
			fn = func(c Context) error {
				return c.String(code, "Unauthorized")
			}
		case http.StatusForbidden:
			fn = func(c Context) error {
				return c.String(code, "Forbidden")
			}
		case http.StatusNotFound:
			fn = func(c Context) error {
				return c.String(code, "Not Found")
			}
		default:
			fn = func(c Context) error {
				return c.String(code, "")
			}
		}
	}
	return fn
}

// Adapter for net/http
func (h Handler) Adapter() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, e error) {
		c := newContext()
		c.Set(CtxRequest, r)
		c.Set(CtxResponse, w)

	}
}
