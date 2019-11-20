package sola

import (
	"net/http"
	"sync"
)

var wg sync.WaitGroup

// Listen & Serve
func Listen(addr string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(addr, g)
	}()
}

// ListenTLS & Serve
func ListenTLS(addr string, certFile string, keyFile string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServeTLS(addr, certFile, keyFile, g)
	}()
}

// Keep Serve
func Keep() {
	wg.Wait()
}
