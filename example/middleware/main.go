package main

// m1 -> Merge(m21 -> m22 -> m23) -> m3

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola/v2"
)

func main() {
	// Sola
	app := sola.New() // 创建 Sola App

	// 过滤
	app.Use(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			if r.URL.String() == "/favicon.ico" {
				return nil
			}
			return next(c)
		}
	})

	// 核心部分
	m2 := sola.Merge(m21, m22, m23)
	app.Use(m1)
	app.Use(m2)
	app.Use(m3)

	// 监听
	sola.Listen("127.0.0.1:3000", app) // 监听 127.0.0.1:3000
	sola.Keep()                        // 固定写法，确保所有监听未结束前程序不退出
}

func m1(next sola.Handler) (handler sola.Handler) {
	fmt.Println("m1 init start")
	handler = func(c sola.Context) (err error) {
		fmt.Println("m1 before")
		err = next(c)
		fmt.Println("m1 after")
		return
	}
	fmt.Println("m1 init end")
	return
}

func m21(next sola.Handler) (handler sola.Handler) {
	fmt.Println("m21 init start")
	handler = func(c sola.Context) (err error) {
		fmt.Println("m21 before")
		err = next(c)
		fmt.Println("m21 after")
		return
	}
	fmt.Println("m21 init end")
	return
}

func m22(next sola.Handler) (handler sola.Handler) {
	fmt.Println("m22 init start")
	handler = func(c sola.Context) (err error) {
		fmt.Println("m22 before")
		err = next(c)
		fmt.Println("m22 after")
		return
	}
	fmt.Println("m22 init end")
	return
}

func m23(next sola.Handler) (handler sola.Handler) {
	fmt.Println("m23 init start")
	handler = func(c sola.Context) (err error) {
		fmt.Println("m23 before")
		err = next(c)
		fmt.Println("m23 after")
		return
	}
	fmt.Println("m23 init end")
	return
}

func m3(next sola.Handler) (handler sola.Handler) {
	fmt.Println("m3 init start")
	handler = func(c sola.Context) (err error) {
		fmt.Println("m3 before")
		// err = next(c)
		err = c.String(http.StatusOK, "Hello World")
		fmt.Println("m3 after")
		return
	}
	fmt.Println("m3 init end")
	return
}
