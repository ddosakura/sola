package proxy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ddosakura/sola/v2"
	lua "github.com/yuin/gopher-lua"
)

// error(s)
var (
	ErrLuaScriptReturn = errors.New("Lua Script should return (int, string)")
	ErrUnsupportCode   = errors.New("Lua Script only support 200, 301")
)

// New Proxy Middleware
func New(script string) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) (e error) {
			r := c.Request()
			w := c.Response()
			L := lua.NewState()
			defer L.Close()
			defer func() {
				if err := recover(); err != nil {
					E, ok := err.(error)
					if ok {
						e = E
					} else {
						e = fmt.Errorf("%v", E)
					}
				}
			}()
			L.SetGlobal("URL", lua.LString(r.URL.String()))
			L.SetGlobal("set_header", L.NewFunction(func(L *lua.LState) int {
				key := L.ToString(1)
				value := L.ToString(2)
				w.Header().Add(key, value)
				return 0
			}))
			if err := L.DoString(script); err != nil {
				return err
			}
			if err := L.CallByParam(lua.P{
				Fn:      L.GetGlobal("handle"),
				NRet:    2,
				Protect: true,
			}); err != nil {
				return err
			}
			ret := L.Get(-2)
			defer L.Pop(2)
			if ret.Type() == lua.LTNil {
				if next == nil {
					return c.Handle(http.StatusNotFound)(c)
				}
				return next(c)
			}
			if ret.Type() != lua.LTNumber {
				return ErrLuaScriptReturn
			}
			code := int(ret.(lua.LNumber))
			data := L.Get(-1).String()
			switch code {
			case http.StatusOK:
				return c.String(code, data)
			case http.StatusMovedPermanently:
				w.Header().Add("Location", data)
				w.WriteHeader(http.StatusMovedPermanently)
				return nil
			}
			return ErrUnsupportCode
		}
	}
}
