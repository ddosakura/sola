package sola

import (
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

// Listen & Serve
func Listen(addr string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := http.ListenAndServe(addr, g); e != nil {
			log.Println(e)
		}
	}()
}

// ListenTLS & Serve
func ListenTLS(addr, certFile, keyFile string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := http.ListenAndServeTLS(addr, certFile, keyFile, g); e != nil {
			log.Println(e)
		}
	}()
}

// Keep Serve
func Keep() {
	wg.Wait()
}
