package auth

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// Clean Data
func Clean(m middleware.Middleware) middleware.Middleware {
	return middleware.Merge(func(c middleware.Context, next middleware.Next) {
		w := c[sola.Response].(http.ResponseWriter)
		http.SetCookie(w, &http.Cookie{Name: authCookieCacheKey})
	}, m)
}
