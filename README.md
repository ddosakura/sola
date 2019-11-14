# sola

A simple golang web framwork based middleware.

## Quick Start

基本的 sola 程序 (Hello World) 如下：

```go
package main

import (
	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

func main() {
	// Sola
	app := sola.New() // 创建 Sola App

	// 核心部分
	app.Use(func(c middleware.Context, next middleware.Next) {
		sola.Text(c, "Hello World") // 输出 Hello World
	})

	// 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep()                        // 固定写法，确保所有监听未结束前程序不退出
}
```

TODO: 框架开发中，更多 demo 请参考 test 目录

## About Writer

Writer 可简化 Response 的书写：

```go
// Text Writer
sola.Text(c, "Hello World")

// JSON Writer
sola.JSON(c, &MyResponse{
	Code: 0,
	Msg:  "Success",
	Data: "Hello World!",
})
```

### Builtin Writer

+ [x] Text	普通文本
+ [x] JSON	JSON 格式

## About Middleware

中间间的定义如下：

```go
// Context for Middleware
type Context map[string]interface{}

// Next func
type Next func()

// Middleware func
type Middleware func(Context, Next)
```

关于 Context 键值的约定：

+ sola      框架
+ router    路由中间件
    + router.param.*    路径参数

### Builtin Middleware

+ [x] auth      认证中间件
+ [x] backup    301 to other host - e.g. http -> https
+ [x] favicon   301 to Online Favicon
+ [x] router    路由中间件
+ [x] static    静态文件中间件

## About ORM

+ [ ] orm
