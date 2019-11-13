package router

import (
	"net/http"
	"strings"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// Meta of Router
type Meta struct {
	method string
	host   string
	urls   []string
	m      middleware.Middleware
}

// Router Middleware
type Router struct {
	routes map[string]*Meta
}

// New Router
func New() *Router {
	return &Router{
		routes: map[string]*Meta{},
	}
}

// Bind "method host/url"
func (r *Router) Bind(match string, m middleware.Middleware) {
	method, host, urls := parse(match)
	r.routes[match] = &Meta{strings.ToUpper(method), host, urls, m}
}

// Routes gen Middleware
func (r *Router) Routes() middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		req := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
	Match:
		for _, meta := range r.routes {
			if meta.method != "" && meta.method != req.Method {
				continue
			}
			if meta.host != "" && meta.host != req.Host {
				continue
			}

			_, _, URLs := parse(req.URL.String())
			if len(URLs) < len(meta.urls) {
				continue
			}
			params := map[string]string{}
			for i, path := range meta.urls {
				if i == 0 { // is host
					continue
				}
				v := URLs[i]
				if path == v {
					continue
				}
				if !strings.HasPrefix(path, ":") {
					continue Match
				}
				k := strings.Replace(path, ":", "", 1)
				params[k] = v
			}
			for k, v := range params {
				c["router.param."+k] = v
			}
			meta.m(c, next)
			return
		}
		if next == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}
		next()
	}
}
