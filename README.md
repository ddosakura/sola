# Sola V2

A simple golang web framwork based middleware.

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

### More Demo

+ [x] [Base Auth](demo/base-auth/main.go)
+ [x] [favicon、static、backup](demo/favicon-static-backup/main.go)
+ [x] [Hello World](demo/hello-world/main.go)
+ [x] [中间件执行顺序](demo/middleware/main.go)
+ [x] [路由&认证 (router、auth)](demo/router-auth/main.go)
+ [x] [完整的注册、登录、鉴权例子](demo/simple-app)

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

+ [x] String	普通文本
+ [x] JSON		JSON 格式

## About Middleware

中间间的定义如下：

```go
type (
	// Context for Middleware
	Context map[string]interface{}

	// Handler func
	Handler func(Context) error

	// Middleware func
	Middleware func(Handler) Handler
)
```

关于 Context 键值的约定：

+ sola      框架
	+ sola.request		*http.Request
	+ sola.response		http.ResponseWriter
+ router    路由中间件
    + router.param.*    路径参数
+ auth		认证中间件
	+ auth.username		Base Auth 用户名
	+ auth.password		Base Auth 密码
	+ auth.claims		JWT Auth Payload

### Builtin Middleware

+ [x] auth      认证中间件
	+ [x] 简化改造
	+ [ ] 自定义返回内容
+ [x] backup    301 to other host - e.g. http -> https
+ [x] favicon   301 to Online Favicon
+ [x] router    路由中间件
	+ [x] 简化改造
+ [x] static    静态文件中间件

## About ORM

+ [x] orm		see gorm
