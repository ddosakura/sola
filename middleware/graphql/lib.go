package graphql

import (
	"github.com/ddosakura/sola/v2"
	lib "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// New GraphQL Handler
func New(s string, resolver interface{}, opts ...lib.SchemaOpt) sola.Handler {
	schema := lib.MustParseSchema(s, resolver, opts...)
	h := &relay.Handler{Schema: schema}
	return func(c sola.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
