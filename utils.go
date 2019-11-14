package sola

import (
	"net/http"

	"github.com/ddosakura/sola/middleware"
)

// SetCookie proxy http.SetCookie
func SetCookie(c middleware.Context, cookie *http.Cookie) {
	w := c[Response].(http.ResponseWriter)
	http.SetCookie(w, cookie)
}
