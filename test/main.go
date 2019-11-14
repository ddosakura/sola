package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/auth"
	"github.com/ddosakura/sola/middleware/backup"
	"github.com/ddosakura/sola/middleware/favicon"
	"github.com/ddosakura/sola/middleware/router"
	"github.com/ddosakura/sola/middleware/static"
)

func hw(c middleware.Context, next middleware.Next) {
	sola.Text(c, "Hello World")
}

func main() {
	r := router.New()
	r.Prefix = "/p" // 设置路由前缀

	// 测试路由匹配
	r.Bind("127.0.0.1:3001/a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A1")
		c[sola.Response].(http.ResponseWriter).Write([]byte("A1"))
	})
	r.Bind("POST /a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A2")
		c[sola.Response].(http.ResponseWriter).Write([]byte("A2"))
	})
	r.Bind("/a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("A3")
		c[sola.Response].(http.ResponseWriter).Write([]byte("A3"))
	})

	// 测试 Merge
	const TMP = "custom.tmp"
	b1 := func(c middleware.Context, next middleware.Next) {
		fmt.Println("B")
		next()
	}
	b2 := func(c middleware.Context, next middleware.Next) {
		next()
		c[sola.Response].(http.ResponseWriter).Write(c[TMP].([]byte))
	}
	b3 := func(c middleware.Context, next middleware.Next) {
		id, e := strconv.Atoi(router.Param(c, "id").(string))
		if e != nil {
			id = -1
		}
		tmp := []byte("UID*2 = " + strconv.Itoa(id*2))
		c[TMP] = tmp
	}
	r.Bind("/b/:id", middleware.Merge(b1, b2, b3))

	r2 := router.New()
	r2.Bind("/b", func(c middleware.Context, next middleware.Next) {
		fmt.Println("r2 - B")
		c[sola.Response].(http.ResponseWriter).Write([]byte("r2 - B"))
	})

	// 测试JWT认证和嵌套路由
	sign := auth.Sign(auth.AuthJWT, []byte("sola_key"))
	AUTH := auth.Auth(auth.AuthJWT, []byte("sola_key"))
	r3 := router.New()
	r31 := router.New()
	r31.Prefix = "/sub"
	r31.Bind("/sub/a", func(c middleware.Context, next middleware.Next) {
		fmt.Println("r31 - 1")
		c[sola.Response].(http.ResponseWriter).Write([]byte("r31 - 1"))
	})
	r31.Bind("/b", func(c middleware.Context, next middleware.Next) {
		fmt.Println("r31 - 2")
		c[sola.Response].(http.ResponseWriter).Write([]byte("r31 - 2"))
	})
	r32 := router.New()
	r32.Prefix = "/sub"
	r32.Bind("/d/:id", func(c middleware.Context, next middleware.Next) {
		claims := c[auth.CtxClaims].(map[string]interface{})
		id := router.Param(c, "id").(string)
		c[sola.Response].(http.ResponseWriter).Write([]byte("No. " + id + "\nUser: " + claims["user"].(string)))
	})
	r3.Bind("/sub", auth.New(AUTH, nil, middleware.Merge(r31.Routes(), r32.Routes())))
	r3.Bind("/login", auth.New(sign, func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		q := r.URL.Query()
		user := q["user"]
		pass := q["pass"]
		if len(user) == 0 || len(pass) == 0 || pass[0] != "123456" {
			c[sola.Response].(http.ResponseWriter).Write([]byte("login fail"))
			return
		}
		c[auth.CtxClaims] = map[string]interface{}{
			"issuer": "sola",
			"user":   user[0],
		}
		next()
	}, func(c middleware.Context, next middleware.Next) {
		c[sola.Response].(http.ResponseWriter).Write([]byte("login success"))
	}))
	r3.Bind("/logout", auth.Clean(func(c middleware.Context, next middleware.Next) {
		c[sola.Response].(http.ResponseWriter).Write([]byte("logout"))
	}))

	// 测试Base认证
	base := auth.Auth(auth.AuthBase, auth.BaseCheck(func(u, p string) bool {
		return u == "admin" && p == "123456"
	}))
	r3.Bind("/base", auth.New(base, nil, hw))

	// 测试 JSON Writer
	r3.Bind("/json", func(c middleware.Context, next middleware.Next) {
		sola.JSON(c, &MyResponse{
			Code: 0,
			Msg:  "Success",
			Data: "Hello World!",
		})
	})

	// 测试 Favicon
	app := sola.New()
	app.Use(favicon.New("http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png"))
	app.Use(r.Routes())
	app.Use(r2.Routes())
	app.Use(r3.Routes())

	// 测试静态文件
	r4 := router.New()
	r4.Prefix = "/static"
	r4.Bind("/", static.New("static"))
	app.Use(r4.Routes())

	// 测试 Backup
	sola.Listen("127.0.0.1:3000", app)
	bak := backup.App("127.0.0.1:3000")
	sola.Listen("127.0.0.1:3001", bak)
	sola.Keep()
}

// MyResponse for json writer
type MyResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
