package main

import (
	"fmt"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/auth"
)

func register(c middleware.Context, next middleware.Next) {
	r := sola.GetRequest(c)
	user := r.PostFormValue("user")
	pass := r.PostFormValue("pass")
	secret := r.PostFormValue("secret")

	if user == "" || pass == "" || secret == "" {
		fail(c)
		return
	}

	if err := db.Create(&User{
		Username: user,
		Password: pass,
		Secret:   secret,
	}).Error; err != nil {
		fail(c)
		return
	}
	next()
}

func login(c middleware.Context, next middleware.Next) {
	r := sola.GetRequest(c)
	user := r.PostFormValue("user")
	pass := r.PostFormValue("pass")

	if user == "" || pass == "" {
		fail(c)
		return
	}

	var u User
	if err := db.First(&u, "username = ?", user).Error; err != nil {
		fail(c)
		return
	}
	if u.Password != pass {
		fail(c)
		return
	}

	fmt.Println(u.ID)
	c[auth.CtxClaims] = map[string]interface{}{
		"id":   u.ID,
		"user": u.Username,
	}
	next()
}
