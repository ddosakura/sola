package sola

import (
	"encoding/json"
	"net/http"
)

// Sola Impl
func (c Context) Sola() *Sola {
	return c[CtxSola].(*Sola)
}

// SetCookie proxy http.SetCookie
func (c Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

// Request in context
func (c Context) Request() *http.Request {
	return c[CtxRequest].(*http.Request)
}

// Response in context
func (c Context) Response() http.ResponseWriter {
	return c[CtxResponse].(http.ResponseWriter)
}

// String Writer
func (c Context) String(statusCode int, data string) error {
	w := c.Response()
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(data))
	return err
}

// JSON Writer
func (c Context) JSON(statusCode int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w := c.Response()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(bs)
	return err
}
