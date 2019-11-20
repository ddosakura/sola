package auth

import (
	"encoding/base64"
	"strings"

	"github.com/ddosakura/sola/v2"
)

func parseBearerAuth(auth string) (token string, ok bool) {
	if len(auth) < len(jwtAuthPrefix) || !strings.EqualFold(auth[:len(jwtAuthPrefix)], jwtAuthPrefix) {
		return
	}
	return auth[len(jwtAuthPrefix):], true
}

func nextFn(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		return next(c)
	}
}

// fork from http pkg
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Claims Reader
func Claims(c sola.Context, key string) interface{} {
	return c[CtxClaims].(map[string]interface{})[key]
}
