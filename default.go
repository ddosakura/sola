package sola

// default
var (
	DefaultApp = New()
)

// Use Middleware
func Use(ms ...Middleware) {
	DefaultApp.Use(ms...)
}
