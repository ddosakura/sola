package hot

import (
	"fmt"

	"github.com/ddosakura/sola/v2"
)

// ctx
const (
	Issuer = "ext/hot"
	CtxHot = "ext.hot"
)

// NotFound Handler
func NotFound(k string) sola.Handler {
	return sola.IError(Issuer, fmt.Errorf("<%s> not found", k)).Handler()
}
