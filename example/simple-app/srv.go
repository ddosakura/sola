package main

import (
	"net/http"
	"strconv"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
)

func welcome(c sola.Context) error {
	user := auth.Claims(c, "user").(string)
	return c.String(http.StatusOK, "Welcome, "+user+"!")
}

func user(c sola.Context) error {
	id, _ := strconv.Atoi(router.Param(c, "id"))
	ID := uint(auth.Claims(c, "id").(float64))

	if uint(id) != ID {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": -1,
			"msg":  "Forbidden",
		})

	}

	var u User
	if err := db.First(&u, ID).Error; err != nil {
		return fail(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": "Your secret: " + u.Secret,
	})
}
