package router

import (
	"net/http"
	"strings"

	"github.com/ddosakura/sola/v2"
)

// Meta of Router
type Meta struct {
	method string
	host   string
	urls   []string
	h      sola.Handler
}

// Router Middleware
type Router struct {
	Prefix string
	routes map[string]*Meta
}

// New Router
func New() *Router {
	return &Router{
		Prefix: "",
		routes: map[string]*Meta{},
	}
}

// Bind "method host/url"
func (r *Router) Bind(match string, m sola.Middleware) {
	method, host, urls := parse(match)
	r.routes[match] = &Meta{strings.ToUpper(method), host, urls, m.Handler()}
}

// BindFunc "method host/url"
func (r *Router) BindFunc(match string, h sola.Handler) {
	method, host, urls := parse(match)
	r.routes[match] = &Meta{strings.ToUpper(method), host, urls, h}
}

// Routes gen Middleware
func (r *Router) Routes() sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			req := c.Request()
		Match:
			for _, meta := range r.routes {
				if meta.method != "" && meta.method != req.Method {
					continue
				}
				if meta.host != "" && meta.host != req.Host {
					continue
				}

				URL := req.URL.EscapedPath()
				if !strings.HasPrefix(URL, r.Prefix) {
					continue
				}
				URL = strings.Replace(URL, r.Prefix, "", 1)
				_, _, URLs := parse(URL)
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
					c[CtxParam(k)] = v
				}
				return meta.h(c)
			}
			if next == nil {
				return c.String(http.StatusNotFound, "Not Found")
			}
			return next(c)
		}
	}
}
