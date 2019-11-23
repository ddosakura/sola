package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ddosakura/sola/v2"
)

// Option of CORS
type Option struct {
	Origin        func(sola.Context) string // Access-Control-Allow-Origin
	ExposeHeaders []string                  // Access-Control-Expose-Headers
	MaxAge        int                       // Access-Control-Max-Age
	Credentials   bool                      // Access-Control-Allow-Credentials
	AllowMethods  []string                  // Access-Control-Allow-Methods
	AllowHeaders  []string                  // Access-Control-Allow-Headers
}

var (
	defaultAllowMethods = []string{
		"GET",
		"PUT",
		"POST",
		"PATCH",
		"DELETE",
		"HEAD",
		"OPTIONS",
	}
)

// New CORS Middleware
// See: https://github.com/zadzbw/koa2-cors/blob/master/src/index.js
func New(o *Option) sola.Middleware {
	if o == nil {
		o = &Option{}
	}
	if o.Origin == nil {
		o.Origin = func(sola.Context) string {
			return "*"
		}
	}
	if o.AllowMethods == nil {
		o.AllowMethods = defaultAllowMethods
	}
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()
			h := w.Header()

			h.Set("Vary", "Origin")

			origin := o.Origin(c)
			if origin == "" {
				return next(c)
			}
			h.Set("Access-Control-Allow-Origin", origin)

			if r.Method == http.MethodOptions {
				if r.Header.Get("Access-Control-Request-Method") == "" {
					return next(c)
				}

				if o.MaxAge >= -1 {
					h.Set("Access-Control-Max-Age", strconv.Itoa(o.MaxAge))
				}
				if o.Credentials {
					h.Set("Access-Control-Allow-Credentials", "true")
				}
				if o.AllowMethods != nil {
					tmp := strings.Join(o.AllowMethods, ",")
					h.Set("Access-Control-Allow-Methods", tmp)
				}
				if o.AllowHeaders == nil {
					tmp := h.Get("Access-Control-Request-Headers")
					h.Set("Access-Control-Allow-Headers", tmp)
				} else {
					tmp := strings.Join(o.AllowHeaders, ",")
					h.Set("Access-Control-Allow-Headers", tmp)
				}

				return c.Handle(http.StatusNoContent)(c)
			}

			if o.Credentials {
				if origin == "*" {
					h.Del("Access-Control-Allow-Credentials")
				} else {
					h.Set("Access-Control-Allow-Credentials", "true")
				}
			}
			if o.ExposeHeaders != nil {
				tmp := strings.Join(o.ExposeHeaders, ",")
				h.Set("Access-Control-Expose-Headers", tmp)
			}

			return next(c)
		}
	}
}
