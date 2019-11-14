package main

import (
	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

func fail(c middleware.Context) {
	sola.JSON(c, map[string]interface{}{
		"code": -1,
		"msg":  "FAIL",
	})
}

func success(c middleware.Context, next middleware.Next) {
	sola.JSON(c, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
	})
}
