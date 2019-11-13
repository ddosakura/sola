package backup

import (
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
)

// Backup Middleware Builder
func Backup(addr string) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
		w.Header().Add("Location", addr+r.URL.String())
		w.WriteHeader(http.StatusMovedPermanently)
	}
}

// App Backup
func App(addr string) *sola.Group {
	app := sola.New()
	app.Use(Backup(addr))
	return app
}
