package auth

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// Clean Data
func Clean(m middleware.Middleware) middleware.Middleware {
	return middleware.Merge(func(c middleware.Context, next middleware.Next) {
		sola.SetCookie(c, &http.Cookie{Name: authCookieCacheKey})
		next()
	}, m)
}
