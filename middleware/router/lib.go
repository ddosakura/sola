package router

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Router Middleware Builder
type Router struct {
	option      *Option
	meta        *Meta
	middlewares []sola.Middleware
	middleware  sola.Middleware
}

// Option of Router
type Option struct {
	Pattern string
	// use sola handler if not match
	UseNotFound bool
}

// New Router Middleware
func New(o *Option) *Router {
	if o == nil {
		o = &Option{}
	}
	return &Router{
		option:      o,
		meta:        buildMeta(o.Pattern),
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
func (r *Router) Sub(o *Option) *Router {
	if o == nil {
		o = &Option{}
	}
	meta := buildMeta(o.Pattern)
	sub := &Router{
		meta:        meta,
		middlewares: []sola.Middleware{},
	}
	m := sola.M(func(c sola.C, next sola.H) error {
		if ctx := match(c, false, meta); ctx != nil {
			NEXT := sola.Handler(next)
			if o.UseNotFound {
				NEXT = c.Handle(http.StatusNotFound)
			}
			return sub.preHandle()(sola.OriginContext(NEXT))(ctx)
		}
		return next(c)
	}).M()
	r.middlewares = append(r.middlewares, m)
	return sub
}

// Use Middleware
func (r *Router) Use(ms ...sola.Middleware) {
	r.middlewares = append(r.middlewares, ms...)
}

// Bind Handler
func (r *Router) Bind(pattern string, h sola.Handler, ms ...sola.Middleware) {
	meta := buildMeta(pattern)
	m := sola.M(func(c sola.C, next sola.H) error {
		if ctx := match(c, true, meta); ctx != nil {
			return sola.MergeFunc(h, ms...)(ctx)
		}
		return next(c)
	}).M()
	r.middlewares = append(r.middlewares, m)
}

// Routes Middleware
func (r *Router) Routes() sola.Middleware {
	m := sola.M(func(c sola.C, next sola.H) error {
		if ctx := match(c, false, r.meta); ctx != nil {
			NEXT := sola.Handler(next)
			if r.option.UseNotFound {
				NEXT = c.Handle(http.StatusNotFound)
			}
			return r.preHandle()(sola.OriginContext(NEXT))(ctx)
		}
		return next(c)
	})
	return m.Must(errNoMatch.Handler())
}
