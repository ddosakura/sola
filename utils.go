package sola

import (
	"net/http"

	"github.com/ddosakura/sola/middleware"
)

// SetCookie proxy http.SetCookie
func SetCookie(c middleware.Context, cookie *http.Cookie) {
	http.SetCookie(GetResponse(c), cookie)
}

// ResponseHeader proxy w.Header()
func ResponseHeader(c middleware.Context) http.Header {
	return GetResponse(c).Header()
}

// GetRequest from context
func GetRequest(c middleware.Context) *http.Request {
	return c[CtxRequest].(*http.Request)
}

// GetResponse from context
func GetResponse(c middleware.Context) http.ResponseWriter {
	return c[CtxResponse].(http.ResponseWriter)
}
