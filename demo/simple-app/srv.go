package main

import (
	"strconv"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/auth"
	"github.com/ddosakura/sola/middleware/router"
)

func welcome(c middleware.Context, next middleware.Next) {
	user := auth.Claims(c, "user").(string)
	sola.Text(c, "Welcome, "+user+"!")
}

func user(c middleware.Context, next middleware.Next) {
	id, _ := strconv.Atoi(router.Param(c, "id"))
	ID := uint(auth.Claims(c, "id").(float64))

	if uint(id) != ID {
		sola.JSON(c, map[string]interface{}{
			"code": -1,
			"msg":  "Forbidden",
		})
		return
	}

	var u User
	if err := db.First(&u, ID).Error; err != nil {
		fail(c)
		return
	}

	sola.JSON(c, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": "Your secret: " + u.Secret,
	})
}
