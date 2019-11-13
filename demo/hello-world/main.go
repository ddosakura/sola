package main

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/favicon"
)

func main() {
	// Sola
	app := sola.New()                                                                                              // 创建 Sola App
	app.Use(favicon.Favicon("http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png")) // 设置 Favicon

	// 核心部分
	app.Use(func(c middleware.Context, next middleware.Next) {
		w := c[sola.Response].(http.ResponseWriter) // 获取 ResponseWriter 对象
		w.Write([]byte("Hello World"))              // 输出 Hello World
	})

	// 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep()                        // 固定写法，确保所有监听未结束前程序不退出
}
