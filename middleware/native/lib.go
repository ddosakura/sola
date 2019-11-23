package native

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Static Middleware
func Static(path, prefix string) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()
			http.StripPrefix(prefix, http.FileServer(http.Dir(path))).ServeHTTP(w, r)
			return nil
		}
	}
}

// From http.Handler
func From(h http.Handler) sola.Handler {
	return func(c sola.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// FromFunc func(http.ResponseWriter, *http.Request)
func FromFunc(h func(http.ResponseWriter, *http.Request)) sola.Handler {
	return func(c sola.Context) error {
		h(c.Response(), c.Request())
		return nil
	}
}
