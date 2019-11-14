package favicon

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// New Favicon Middleware
func New(url string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
		if r.URL.String() == "/favicon.ico" {
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusMovedPermanently)
		} else if next != nil {
			next()
		}
	}
}
