package linux

import (
	"log"
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
		if e := endless.ListenAndServe(addr, g); e != nil {
			log.Println(e)
		}
	}()
}

// ListenTLS & Serve
func ListenTLS(addr, cert, key string, g *sola.Sola) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := endless.ListenAndServeTLS(addr, cert, key, g); e != nil {
			log.Println(e)
		}
	}()
}

// Keep Serve
func Keep() {
	wg.Wait()
}
