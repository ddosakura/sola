package main

import (
	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/auth"
)

// 用户名、密码验证函数
var check = auth.BaseCheck(func(u, p string) bool {
	return u == "admin" && p == "123456"
})

func main() {
	// Sola
	app := sola.New()                       // 创建 Sola App
	base := auth.Auth(auth.AuthBase, check) // 创建 Base Auth

	app.Use(auth.New(base, nil, func(c middleware.Context, next middleware.Next) {
		sola.Text(c, "Hello World") // 输出 Hello World
	}))

	// 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep()                        // 固定写法，确保所有监听未结束前程序不退出
}
