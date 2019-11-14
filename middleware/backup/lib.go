package backup

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// New Backup Middleware
func New(addr string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := sola.GetRequest(c)
		sola.ResponseHeader(c).Add("Location", addr+r.URL.String())
		sola.GetResponse(c).WriteHeader(http.StatusMovedPermanently)
	}
}

// App Backup
func App(addr string) *sola.Group {
	app := sola.New()
	app.Use(New(addr))
	return app
}
