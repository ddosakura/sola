package middleware

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
