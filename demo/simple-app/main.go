package main

import (
	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
)

func main() {
	defer db.Close()
	app := sola.New()
	r := router.New()

	r.BindFunc("POST /register", register)
	r.BindFunc("POST /login", auth.NewFunc(_sign, login, success))
	r.BindFunc("/logout", auth.CleanFunc(success))
	r.BindFunc("GET /welcome", auth.NewFunc(_auth, nil, welcome))
	r.BindFunc("GET /user/:id", auth.NewFunc(_auth, nil, user))

	app.Use(r.Routes())
	sola.Listen("127.0.0.1:3000", app)
	sola.Keep()
}
