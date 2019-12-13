package proxy

import (
	"errors"

	"github.com/ddosakura/sola/v2"
)

// Ctx
const (
	Issuer = "sola/proxy"
)

// error(s)
var (
	ErrLuaScriptReturn = errors.New("Lua Script should return (int, string)")
	ErrUnsupportCode   = errors.New("Lua Script only support 200, 301")

	errLuaScriptReturn = sola.IError(Issuer, ErrLuaScriptReturn)
	errUnsupportCode   = sola.IError(Issuer, ErrUnsupportCode)
)
