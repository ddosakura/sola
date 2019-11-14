package main

import (
	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware/backup"
	"github.com/ddosakura/sola/middleware/favicon"
	"github.com/ddosakura/sola/middleware/router"
	"github.com/ddosakura/sola/middleware/static"
)

func main() {
	app := sola.New()

	// 设置 favicon
	app.Use(favicon.New("http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png"))

	// 和路由中间件一起使用
	r := router.New()
	r.Prefix = "/s"
	r.Bind("", static.New("static", "/s"))
	app.Use(r.Routes())
	// 直接使用
	app.Use(static.New(".", ""))

	sola.Listen("127.0.0.1:3000", app)

	// 自动跳转 http://127.0.0.1:3000
	bak := backup.App("http://127.0.0.1:3000")
	sola.Listen("127.0.0.1:3001", bak)

	sola.Keep()
}
