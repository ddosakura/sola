package sola

import (
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func listen(addr string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := http.ListenAndServe(addr, g); e != nil {
			log.Println(e)
		}
	}()
}

// Listen & Serve(s)
func Listen(addr string, gs ...*Sola) {
	if len(gs) == 0 {
		listen(addr, DefaultApp)
		return
	}
	for _, g := range gs {
		listen(addr, g)
	}
}

func listenTLS(addr, certFile, keyFile string, g *Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := http.ListenAndServeTLS(addr, certFile, keyFile, g); e != nil {
			log.Println(e)
		}
	}()
}

// ListenTLS & Serve(s)
func ListenTLS(addr, certFile, keyFile string, gs ...*Sola) {
	if len(gs) == 0 {
		listenTLS(addr, certFile, keyFile, DefaultApp)
		return
	}
	for _, g := range gs {
		listenTLS(addr, certFile, keyFile, g)
	}
}

// Keep Serve(s)
func Keep() {
	wg.Wait()
}

// ListenKeep Serve(s)
func ListenKeep(addr string, gs ...*Sola) {
	Listen(addr, gs...)
	Keep()
}

// ListenKeepTLS Serve(s)
func ListenKeepTLS(addr, certFile, keyFile string, gs ...*Sola) {
	ListenTLS(addr, certFile, keyFile, gs...)
	Keep()
}
