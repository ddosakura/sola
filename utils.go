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

// ResponseHeader proxy w.Header()
func ResponseHeader(c middleware.Context) http.Header {
	w := c[Response].(http.ResponseWriter)
	return w.Header()
}

// GetRequest from context
func GetRequest(c middleware.Context) *http.Request {
	return c[Request].(*http.Request)
}
