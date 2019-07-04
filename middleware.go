package middlewares

import "net/http"

// Middleware HTTP middleware
type Middleware func(http.Handler) http.Handler

// Chain Chain of middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next
	}
}
