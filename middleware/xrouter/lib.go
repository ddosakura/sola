package xrouter

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Router Middleware Builder
type Router struct {
	meta *Meta

	middlewares []sola.Middleware
	middleware  sola.Middleware
}

// New Router Middleware
func New(pattern string) *Router {
	return &Router{
		meta:        buildMeta(pattern),
		middlewares: []sola.Middleware{},
	}
}

func (r *Router) preHandle() sola.Middleware {
	if r.middleware == nil {
		r.middleware = sola.Merge(r.middlewares...)
	}
	return r.middleware
}

// Sub Router
func (r *Router) Sub(pattern string) *Router {
	meta := buildMeta(pattern)
	sub := &Router{
		meta:        meta,
		middlewares: []sola.Middleware{},
	}
	fn := func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			if quick := match(c, false, meta); quick != nil {
				quick()
				return sub.preHandle()(c.Handle(http.StatusNotFound))(c)
			}
			return next(c)
		}
	}
	r.middlewares = append(r.middlewares, fn)
	return sub
}

// Use Middleware
func (r *Router) Use(m sola.Middleware) {
	r.middlewares = append(r.middlewares, m)
}

// Bind Handler
func (r *Router) Bind(pattern string, h sola.Handler) {
	meta := buildMeta(pattern)
	fn := func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			if quick := match(c, true, meta); quick != nil {
				return h(c)
			}
			return next(c)
		}
	}
	r.middlewares = append(r.middlewares, fn)
}

// Routes Middleware
func (r *Router) Routes() sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			if quick := match(c, false, r.meta); quick != nil {
				quick()
				return r.preHandle()(c.Handle(http.StatusNotFound))(c)
			}
			return next(c)
		}
	}
}
