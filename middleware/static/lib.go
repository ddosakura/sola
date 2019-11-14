package static

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// New Static Middleware
func New(path, prefix string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := sola.GetRequest(c)
		w := sola.GetResponse(c)
		http.StripPrefix(prefix, http.FileServer(http.Dir(path))).ServeHTTP(w, r)
	}
}
