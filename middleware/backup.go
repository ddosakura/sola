package middleware

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Backup Middleware
func Backup(addr string) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()
			w.Header().Add("Location", addr+r.URL.String())
			w.WriteHeader(http.StatusMovedPermanently)
			return nil
		}
	}
}

// BackupSola App
func BackupSola(addr string) *sola.Sola {
	app := sola.New()
	app.Use(Backup(addr))
	return app
}
