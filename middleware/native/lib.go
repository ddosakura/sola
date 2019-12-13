package native

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Static Middleware
func Static(path, prefix string) sola.Middleware {
	h := http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
	return sola.FromHandler(h).M()
}
