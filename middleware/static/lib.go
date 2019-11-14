package static

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// New Static Middleware
func New(path, prefix string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
		http.StripPrefix(prefix, http.FileServer(http.Dir(path))).ServeHTTP(w, r)
	}
}
