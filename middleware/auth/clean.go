package auth

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Clean Data
func Clean(m sola.Middleware) sola.Middleware {
	return sola.Merge(sola.Handler(func(c sola.Context) error {
		c.SetCookie(&http.Cookie{Name: authCookieCacheKey})
		return nil
	}).M(), m)
}

// CleanFunc Data
func CleanFunc(h sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		c.SetCookie(&http.Cookie{Name: authCookieCacheKey})
		return h(c)
	}
}
