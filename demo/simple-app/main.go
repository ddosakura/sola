package main

import (
	"github.com/ddosakura/sola/middleware"

	// 按需加载
	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware/auth"
	"github.com/ddosakura/sola/middleware/router"
)

func main() {
	defer db.Close()
	app := sola.New()
	r := router.New()

	r.Bind("POST /register", middleware.Merge(register, success))
	r.Bind("POST /login", auth.New(_sign, login, success))
	r.Bind("/logout", auth.Clean(success))
	r.Bind("GET /welcome", auth.New(_auth, nil, welcome))
	r.Bind("GET /user/:id", auth.New(_auth, nil, user))

	app.Use(r.Routes())
	sola.Listen("127.0.0.1:3000", app)
	sola.Keep()
}
