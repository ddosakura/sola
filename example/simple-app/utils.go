package main

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

var fail = sola.Handler(func(c sola.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": -1,
		"msg":  "FAIL",
	})
})

var success = sola.Handler(func(c sola.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
	})
})
