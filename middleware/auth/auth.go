package auth

import (
	"errors"

	"github.com/ddosakura/sola/v2"
)

// Type of Auth
type Type uint8

// AuthType
const (
	AuthBase Type = iota
	AuthJWT
)

// BaseCheck for Base Auth
type BaseCheck func(username, password string) bool

func basePass(string, string) bool {
	return true
}

const jwtAuthPrefix = "Bearer "
const authCookieCacheKey = "Sola-Authorization"

// Context Key
const (
	Issuer      = "sola/auth"
	CtxUsername = "auth.username"
	CtxPassword = "auth.password"
	CtxClaims   = "auth.claims"
	CtxToken    = "auth.token"
)

// error(s)
var (
	ErrNoClaims = errors.New("auth.claims not found")

	errNoClaims = sola.IError(Issuer, ErrNoClaims)
)
