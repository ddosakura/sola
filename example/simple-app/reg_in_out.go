package main

import (
	"fmt"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
)

func register(c sola.Context) error {
	r := c.Request()
	user := r.PostFormValue("user")
	pass := r.PostFormValue("pass")
	secret := r.PostFormValue("secret")

	if user == "" || pass == "" || secret == "" {
		return fail(c)
	}

	if err := db.Create(&User{
		Username: user,
		Password: pass,
		Secret:   secret,
	}).Error; err != nil {
		return fail(c)
	}
	return success(c)
}

func login(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		r := c.Request()
		user := r.PostFormValue("user")
		pass := r.PostFormValue("pass")

		if user == "" || pass == "" {
			return fail(c)
		}

		var u User
		if err := db.First(&u, "username = ?", user).Error; err != nil {
			return fail(c)
		}
		if u.Password != pass {
			return fail(c)
		}

		fmt.Println(u.ID)
		c[auth.CtxClaims] = map[string]interface{}{
			"id":   u.ID,
			"user": u.Username,
		}
		return next(c)
	}
}
