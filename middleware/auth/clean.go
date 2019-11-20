package auth

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Clean Data
func Clean(m sola.Middleware) sola.Middleware {
	return sola.Merge(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c.SetCookie(&http.Cookie{Name: authCookieCacheKey})
			return next(c)
		}
	}, m)
}

// CleanFunc Data
func CleanFunc(h sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		c.SetCookie(&http.Cookie{Name: authCookieCacheKey})
		return h(c)
	}
}
