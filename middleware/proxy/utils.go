package proxy

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ddosakura/sola/v2"
)

// Balance Proxy
func Balance(c sola.Context, targets []*url.URL) *httputil.ReverseProxy {
	if len(targets) < 1 {
		c.Handle(http.StatusNotFound)
		return nil
	}
	target := targets[rand.Int()%len(targets)]
	if target == nil {
		c.Handle(http.StatusNotFound)
		return nil
	}
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}
