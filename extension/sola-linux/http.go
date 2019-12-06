package linux

import (
	"sync"

	"github.com/ddosakura/sola/v2"
	"github.com/fvbock/endless"
)

var wg sync.WaitGroup

// Listen & Serve
func Listen(addr string, g *sola.Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		endless.ListenAndServe(addr, g)
	}()
}

// ListenTLS & Serve
func ListenTLS(addr string, certFile string, keyFile string, g *sola.Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		endless.ListenAndServeTLS(addr, certFile, keyFile, g)
	}()
}

// Keep Serve
func Keep() {
	wg.Wait()
}
