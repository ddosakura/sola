package middleware

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Favicon Middleware
func Favicon(url string) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()
			if r.URL.String() == "/favicon.ico" {
				w.Header().Add("Location", url)
				w.WriteHeader(http.StatusMovedPermanently)
				return nil
			}

			if next != nil {
				return next(c)
			}
			return nil
		}
	}
}
