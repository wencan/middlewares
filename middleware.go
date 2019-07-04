package middlewares

import (
	"net/http"
)

// Middleware HTTP middleware
type Middleware func(http.Handler) http.Handler

// Chain Chain of middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}
		return next
	}
}
