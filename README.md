# Sola V2

A simple golang web framwork based middleware.

+ [Change Log](./CHANGELOG.md)

## 内部错误（无法修复）

+ v2.1.0 & v2.1.1 内部版本号变量误标记为 `2.0.0`

## Quick Start

基本的 sola 程序 (Hello World) 如下：

```go
package main

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// handler
func hw(c sola.Context) error {
	// 输出 Hello World
	return c.String(http.StatusOK, "Hello World")
}

func main() {
	// 将 handler 包装为中间件
	m := sola.Handler(hw).M()
	// 使用中间件
	sola.Use(m)
	// 监听 127.0.0.1:3000
	sola.ListenKeep("127.0.0.1:3000")
}
```

## Hot Update

```bash
# 安装
go get -u -v github.com/ddosakura/sola/v2/cli/sola-hot
# 在开发目录执行
sola-hot
```

执行 `sola-hot` 进行热更新，运行过程中将生成临时可执行文件 `sola-dev`。

linux 系统使用 `sola-linux` 监听端口可实现平滑切换（不中断请求）。 

```go
import (
	linux "github.com/ddosakura/sola/v2/extension/sola-linux"
)

linux.Listen("127.0.0.1:3000", app)
linux.Keep()
```

## Extension

+ hot			动态模块加载（使用了 plugin，仅 linux 可用）
+ sola-linux 	平滑切换（使用了大量系统调用，仅 linux 可用）

### More Example

+ [Example 仓库地址](https://github.com/it-repo/box-example)

## 废弃说明

即将废弃的中间件/方法会在注释中标注 `@deprecated`

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
		Store() map[string]interface{}
		Origin() Context
		Shadow() Context
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
		origin Context
		lock   sync.RWMutex
		store  map[string]interface{}
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
+ auth		认证中间件
	+ auth.username		Base Auth 用户名
	+ auth.password		Base Auth 密码
	+ auth.claims		JWT Auth Payload
	+ auth.token        签发的 JWT
+ logger	日志中间件
	+ logger			message chan
+ x/router  旧路由中间件
    + x.router.param.*  路径参数
+ router    新路由中间件
	+ router.meta		路由元数据
    + router.param.*    路径参数

+ hot 		动态模块加载扩展
	+ ext.hot			*Hot

> **注意：请尽量使用中间件包中提供的标准方法、常量读取上下文数据，防止内部逻辑修改导致的向后不兼容。**

### Builtin Middleware

+ [x] auth      认证中间件
	+ [x] 自定义返回内容
	+ [x] Dev Mode(500)
+ [x] cors		跨域中间件 - 参考 [koa2-cors](https://github.com/zadzbw/koa2-cors)
+ [x] graphql   GraphQL 中间件
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
	+ [x] 负载均衡
+ [x] rest      RESTful API 中间件
+ [x] swagger   API 文档中间件
+ [x] ws		WebSocket 中间件
+ [x] x/router  旧路由中间件 (@deprecated 预计在 v2.2.x 移除)
+ [x] router	新路由中间件

#### router 匹配规则

+ 完整版 `GET localhost:3000/user/:id` (注意 Host 直接匹配 *http.Request 中的 Host)
+ 不指定 Method `localhost:3000/user/:id`
+ 不指定 Host `GET /user/:id`
+ 不指定 Method & Host `/user/:id`
+ 默认路径匹配（仅在 Bind 中使用，一般作为最后一个路径） `/*`
+ 特殊用法：仅指定 Method（用于 Bind，Path 由 Router 的 Pattern 匹配） `GET` - 支持 `net/http` 包中的所有 `Method*`

样例解释：

```go
// 匹配 /api/v1/*
r := router.New(&router.Option{
	Pattern: "/api/v1",
})
// 严格匹配 /api/v1/hello
r.Bind("/hello", hello("Hello World!"))
{
	// 匹配 /api/v1/user/*
	sub := r.Sub(&router.Option{
		Pattern: "/user",
	})
	// 严格匹配 /api/v1/user/hello
	sub.Bind("/hello", hello("Hello!"))
	// 严格匹配 /api/v1/user/hello 失败后执行该中间件
	sub.Use(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			fmt.Println("do auth")
			return next(c)
		}
	})
	// 严格匹配 /api/v1/user/info
	sub.Bind("/info", hello("user info"))
	// 匹配 /api/v1/user/infox/*
	sub.Bind("/infox/*", hello("user infox"))
	// 严格匹配 /api/v1/user/:id （路径参数 id）
	sub.Bind("/:id", get("id"))
	// sub 没有默认匹配
	// 如果 sub 加了 UseNotFound 选项，将调用 sola 的 404 Handler
	// 否则不会调用，将继续匹配后面的 `/*`
}
// 匹配 /api/v1/* (默认匹配)
r.Bind("/*", error404)
// 使用路由
app.Use(r.Routes())
```

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
	"github.com/ddosakura/sola/v2/middleware/x/router"
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
