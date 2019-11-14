package sola

import (
	"encoding/json"
	"net/http"

	"github.com/ddosakura/sola/middleware"
)

// JSON Writer
func JSON(c middleware.Context, data interface{}, statusCode ...int) (int, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	w := c[Response].(http.ResponseWriter)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	return w.Write(bs)
}

// Text Writer
func Text(c middleware.Context, data string, statusCode ...int) (int, error) {
	w := c[Response].(http.ResponseWriter)
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	return w.Write([]byte(data))
}
