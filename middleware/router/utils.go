package router

import (
	"strings"

	"github.com/ddosakura/sola/middleware"
)

func parse(match string) (method string, host string, urls []string) {
	matches := strings.Split(match, " ")
	url := matches[0]
	if len(matches) > 1 {
		method = matches[0]
		url = matches[1]
	}
	urls = strings.Split(url, "/")
	host = urls[0]
	return
}

// CtxParam Builder
func CtxParam(key string) string {
	return "router.param." + key
}

// Param in route
func Param(c middleware.Context, key string) string {
	return c[CtxParam(key)].(string)
}
