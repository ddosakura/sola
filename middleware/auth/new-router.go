package auth

import (
	"github.com/ddosakura/sola/v2"
)

// NewJWT Middleware
func NewJWT(key []byte) (sign sola.Middleware, wrapper sola.Middleware) {
	sign = Sign(AuthJWT, key)
	auth := Auth(AuthJWT, key)
	return sign, sola.Merge(loadAuthCache, auth)
}

// LoadAuthCache Middleware's export
var LoadAuthCache = loadAuthCache
