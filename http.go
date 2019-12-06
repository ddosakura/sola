package sola

import (
	"sync"

	"github.com/fvbock/endless"
)

var wg sync.WaitGroup

// Listen & Serve
func Listen(addr string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		endless.ListenAndServe(addr, g)
	}()
}

// ListenTLS & Serve
func ListenTLS(addr string, certFile string, keyFile string, g *Sola) {
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
