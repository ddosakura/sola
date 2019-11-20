package main

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

func main() {
	// Sola
	app := sola.New() // 创建 Sola App

	// 核心部分
	app.Use(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			// 输出 Hello World
			return c.String(http.StatusOK, "Hello World")
		}
	})

	// 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep()                        // 固定写法，确保所有监听未结束前程序不退出
}
