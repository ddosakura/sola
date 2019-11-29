# Sola V2

A simple golang web framwork based middleware.

+ [Change Log](./CHANGELOG.md)

## Quick Start

基本的 sola 程序 (Hello World) 如下：

```go
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
```

### More Example

+ [Example 仓库地址](https://github.com/it-repo/box-example)

## About Reader

Reader 可简化 Request 的读取：

```go
var a ReqLogin
if err := c.GetJSON(&a); err != nil {
	return err
}
```

### Builtin Reader

+ [x] GetJSON   JSON

## About Writer

Writer 可简化 Response 的书写：

```go
// String Writer
c.String(http.StatusOK, "Hello World")

// JSON Writer
c.JSON(http.StatusOK, &MyResponse{
	Code: 0,
	Msg:  "Success",
	Data: "Hello World!",
})
```

### Builtin Writer

+ [x] Blob		二进制
+ [x] HTML      HTML(text/html)
+ [x] String	普通文本(text/plain)
+ [x] JSON		JSON(application/json)
+ [x] File      文件 - 兼容 afero

## About Middleware

中间间的定义如下：

```go
type (
	// Context for Middleware
	Context interface {
		// Set/Get
		Set(key string, value interface{})
		Get(key string) interface{}

		// API
		Sola() *Sola
		SetCookie(cookie *http.Cookie)
		Request() *http.Request
		Response() http.ResponseWriter

		// Writer
		Blob(code int, contentType string, bs []byte) (err error)
		HTML(code int, data string) error
		String(code int, data string) error
		JSON(code int, data interface{}) error
		File(f File) (err error)

		// Reader
		GetJSON(data interface{}) error

		// Handler
		Handle(code int) Handler
	}
	context struct {
		lock  sync.RWMutex
		store map[string]interface{}
	}

	// Handler func
	Handler func(Context) error

	// Middleware func
	Middleware func(Handler) Handler
)
```

关于 Context 键值的约定：

+ sola      框架
	+ sola				*sola.Sola
	+ sola.request		*http.Request
	+ sola.response		http.ResponseWriter
+ router    路由中间件
    + router.param.*    路径参数
+ auth		认证中间件
	+ auth.username		Base Auth 用户名
	+ auth.password		Base Auth 密码
	+ auth.claims		JWT Auth Payload
	+ auth.token        签发的 JWT
+ logger	日志中间件
	+ logger			message chan

### Builtin Middleware

+ [x] auth      认证中间件
	+ [x] 自定义返回内容
	+ [x] Dev Mode(500)
+ [x] cors		跨域中间件 - 参考 [koa2-cors](https://github.com/zadzbw/koa2-cors)
+ [x] logger    日志中间件
+ [x] native	go 原生 handler 转换中间件(取代原静态文件中间件)
	+ [x] static    原静态文件中间件
	+ 可用于静态文件
	+ 可用于 statik
	+ 可用于 afero
	+ ...
+ [x] proxy     反向代理中间件(取代原 backup、favicon 中间件)
	+ [x] backup    301 to other host - e.g. http -> https
	+ [x] favicon   301 to Online Favicon
	+ 嵌入 lua 脚本：https://github.com/yuin/gopher-lua
	+ [ ] 完善
+ [x] rest      RESTful API 中间件
+ [x] router    路由中间件
	+ [ ] 中间件改造 xrouter?
+ [ ] swagger   API 文档中间件

## About Config

+ [x] Dev Mode
+ see [viper](https://github.com/spf13/viper)

## About ORM

+ [x] Debug in Dev Mode
+ see [gorm](https://github.com/jinzhu/gorm)

## About API Doc

+ [swagger](https://github.com/swaggo/swag)

以下例子仅作为 API Doc 使用说明，未实现具体功能：

```go
package main

import (
	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/ddosakura/sola/v2/middleware/swagger"

	_ "example/sola-example/api-doc/docs"
	"example/sola-example/api-doc/handler"
)

// @title Swagger Example API
// @version 1.0
// @host localhost:3000
// @BasePath /api/v1
// @description This is a sample server celler server.

// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	_sign := auth.Sign(auth.AuthJWT, []byte("sola_key"))
	_auth := auth.Auth(auth.AuthJWT, []byte("sola_key"))

	app := sola.New()
	r := router.New()

	r.BindFunc("GET /swagger", swagger.WrapHandler)

	sub := router.New()
	sub.Prefix = "/api/v1"
	{
		sub.BindFunc("GET /hello", handler.Hello)
		sub.BindFunc("POST /login", auth.NewFunc(_sign, tmp, handler.Hello))
		sub.BindFunc("/logout", auth.CleanFunc(handler.Hello))

		third := router.New()
		third.Prefix = sub.Prefix
		{
			third.BindFunc("GET /list", handler.List)
			third.BindFunc("GET /item/:id", handler.Item)
		}
		sub.Bind("", auth.New(_auth, nil, third.Routes()))
	}
	r.Bind("/api/v1", sub.Routes())

	app.Use(r.Routes())
	sola.Listen("127.0.0.1:3000", app)
	sola.Keep()
}

func tmp(sola.Handler) sola.Handler {
	return handler.Hello
}
```

```go
// Hello godoc
// @Summary     Say Hello
// @Description Print Hello World!
// @Produce     plain
// @Success     200 {string} string "Hello World!"
// @Router      /hello [get]
func Hello(c sola.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}
```
