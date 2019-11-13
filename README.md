# sola

A simple golang web framwork based middleware.

## Quick Start

基本的 sola 程序 (Hello World) 如下：

```go
package main

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/ddosakura/sola/middleware/favicon"
)

func main() {
    // Sola
	app := sola.New() // 创建 Sola App
    app.Use(favicon.Favicon("http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png")) // 设置 Favicon

	// 核心部分
	app.Use(func(c middleware.Context, next middleware.Next) {
		w := c[sola.Response].(http.ResponseWriter) // 获取 ResponseWriter 对象
		w.Write([]byte("Hello World")) // 输出 Hello World
	})

    // 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep() // 固定写法，确保所有监听未结束前程序不退出
}
```

TODO: 框架开发中，更多 demo 请参考 test 目录

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

+ [ ] auth      认证中间件
+ [x] backup    301 to other host - e.g. http -> https
+ [x] favicon   301 to Online Favicon
+ [x] router    路由中间件

## About ORM

+ [ ] orm
