package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/backup"
	"github.com/ddosakura/sola/middleware/favicon"
	"github.com/ddosakura/sola/middleware/router"
)

func hw(c middleware.Context, next middleware.Next) {
	w := c[sola.Response].(http.ResponseWriter)
	w.Write([]byte("Hello World!"))
}

func main() {
	r := router.New()

	// 测试路由匹配
	r.Bind("127.0.0.1:3001/a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A1")
	})
	r.Bind("POST /a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A2")
	})
	r.Bind("/a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A3")
	})

	// 测试 Merge
	b1 := func(c middleware.Context, next middleware.Next) {
		fmt.Println("B")
		next()
	}
	b2 := func(c middleware.Context, next middleware.Next) {
		next()
		c[sola.Response].(http.ResponseWriter).Write(c["tmp"].([]byte))
	}
	b3 := func(c middleware.Context, next middleware.Next) {
		id, e := strconv.Atoi(c["router.param.id"].(string))
		if e != nil {
			id = -1
		}
		tmp := []byte("UID*2 = " + strconv.Itoa(id*2))
		c["tmp"] = tmp
	}
	r.Bind("/b/:id", middleware.Merge(b1, b2, b3))

	// 测试 Favicon
	app := sola.New()
	app.Use(favicon.Favicon("http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png"))
	app.Use(r.Routes())

	// 测试 Backup
	sola.Listen("127.0.0.1:3000", app)
	bak := backup.App("127.0.0.1:3000")
	sola.Listen("127.0.0.1:3001", bak)
	sola.Keep()
}
