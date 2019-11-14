package sola

import (
	"encoding/json"

	"github.com/ddosakura/sola/middleware"
)

// JSON Writer
func JSON(c middleware.Context, data interface{}, statusCode ...int) (int, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	w := GetResponse(c)
	ResponseHeader(c).Set("Content-Type", "application/json; charset=utf-8")
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	return w.Write(bs)
}

// Text Writer
func Text(c middleware.Context, data string, statusCode ...int) (int, error) {
	w := GetResponse(c)
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	return w.Write([]byte(data))
}
