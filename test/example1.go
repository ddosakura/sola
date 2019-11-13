package main

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

var key = "text"

// 基本中间件框架测试
func example1() {
	app := sola.New()

	app.Use(func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
		if r.URL.String() != "/favicon.ico" {
			next()
			s := c[key].(string)
			w.Write([]byte(s))
		}
	})

	app.Use(func(c middleware.Context, next middleware.Next) {
		fmt.Println("M1 Start", c)
		next()
		s := c[key].(string) + "!"
		c[key] = s
		fmt.Println("M1 Finish", c)
	})

	app.Use(func(c middleware.Context, next middleware.Next) {
		fmt.Println("M2 Start", c)
		next()
		s := c[key].(string) + " World"
		c[key] = s
		fmt.Println("M2 Finish", c)
	})

	app.Use(func(c middleware.Context, next middleware.Next) {
		fmt.Println("M3 Start", c)
		c[key] = "Hello"
		fmt.Println("M3 Finish", c)
	})

	http.ListenAndServe("127.0.0.1:3000", app)
}
