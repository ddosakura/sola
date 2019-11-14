package favicon

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// New Favicon Middleware
func New(url string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := sola.GetRequest(c)
		if r.URL.String() == "/favicon.ico" {
			sola.ResponseHeader(c).Add("Location", url)
			sola.GetResponse(c).WriteHeader(http.StatusMovedPermanently)
			return
		}

		if next != nil {
			next()
		}
	}
}
