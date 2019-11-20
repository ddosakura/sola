package main

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
)

func main() {
	app := sola.New()

	// 定义秘钥
	sign := auth.Sign(auth.AuthJWT, []byte("sola_key"))
	AUTH := auth.Auth(auth.AuthJWT, []byte("sola_key"))

	sub := router.New()
	sub.Prefix = "/user" // 路由前缀
	sub.BindFunc("/:id", func(c sola.Context) error {
		id := router.Param(c, "id") // 路径参数

		// 获取用户名
		claims := c[auth.CtxClaims].(map[string]interface{})
		user := claims["user"].(string)

		// 输出
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
			"data": fmt.Sprintf("ID = %s, USER = %s", id, user),
		})
	})

	r := router.New()
	r.Bind("/user", auth.New(AUTH, nil, sub.Routes())) // 二级路由 /user 需要登录

	// 登录（必须使用 GET 方式）
	r.BindFunc("GET /login", auth.NewFunc(sign, func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			// 获取 GET 参数
			r := c.Request()
			q := r.URL.Query()
			user := q["user"]
			pass := q["pass"]

			// 校验
			if len(user) == 0 || len(pass) == 0 || pass[0] != "123456" {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"code": -1,
					"msg":  "FAIL",
				})
			}

			// 储存用户名等信息
			c[auth.CtxClaims] = map[string]interface{}{
				"issuer": "sola",
				"user":   user[0],
			}
			return next(c) // 登录成功
		}
	}, func(c sola.Context) error {
		// 登录成功调用
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
		})
	}))
	// 清除登录状态（必须使用 GET 方式）
	r.BindFunc("GET /logout", auth.CleanFunc(func(c sola.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
		})
	}))

	app.Use(r.Routes()) // 一级路由
	sola.Listen("127.0.0.1:3000", app)
	sola.Keep()
}
