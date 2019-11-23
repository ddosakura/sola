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

### More Example

+ [Example 仓库地址](https://github.com/it-repo/box-example)

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
	+ sola				*sola.Sola
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
	+ [x] 自定义返回内容
	+ [x] Dev Mode(500)
+ [x] cors		跨域中间件 - 参考 [koa2-cors](https://github.com/zadzbw/koa2-cors)
+ [x] native	go 原生 handler 转换中间件(取代原静态文件中间件)
	+ [x] static    原静态文件中间件
	+ 可用于静态文件
	+ 可用于 statik
	+ 可用于 afero
	+ ...
+ [ ] proxy     反向代理中间件(取代原 backup、favicon 中间件)
	+ 考虑嵌入 lua 脚本：https://github.com/yuin/gopher-lua
+ [x] router    路由中间件
+ [x] backup    301 to other host - e.g. http -> https
+ [x] favicon   301 to Online Favicon

## About Config

+ [x] Dev Mode
+ see [viper](https://github.com/spf13/viper)

## About ORM

+ [x] Debug in Dev Mode
+ see [gorm](https://github.com/jinzhu/gorm)
