package sola

import (
	"net/http"

	"github.com/ddosakura/sola/middleware"
)

// Group of Middleware
type Group struct {
	middlewares []middleware.Middleware
}

// New Middleware Group
func New() *Group {
	return &Group{
		middlewares: []middleware.Middleware{},
	}
}

// Use Middleware
func (g *Group) Use(m middleware.Middleware) {
	g.middlewares = append(g.middlewares, m)
}

// ServeHTTP to impl http handler
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := middleware.Context{}
	c[Request] = r
	c[Response] = w

	middleware.Merge(g.middlewares...)(c, nil)
}
