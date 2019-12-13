package sola

import (
	"net/http"
)

// SetHandler in Sola
func (s *Sola) SetHandler(code int, h Handler) {
	s.handlers[code] = h
}

// Custom Handle Code
const (
	HandleCodePass = iota
)

// HandlePass Handler
var HandlePass = Handler(func(c Context) error {
	return nil
})

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
