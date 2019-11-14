package auth

import (
	"encoding/base64"
	"strings"

	"github.com/ddosakura/sola/middleware"
)

func parseBearerAuth(auth string) (token string, ok bool) {
	if len(auth) < len(jwtAuthPrefix) || !strings.EqualFold(auth[:len(jwtAuthPrefix)], jwtAuthPrefix) {
		return
	}
	return auth[len(jwtAuthPrefix):], true
}

func nextFn(c middleware.Context, next middleware.Next) {
	next()
}

// fork from http pkg
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Claims Reader
func Claims(c middleware.Context, key string) interface{} {
	return c[CtxClaims].(map[string]interface{})[key]
}
