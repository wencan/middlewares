package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryHandlerFunc Handler function that handle recovery
type RecoveryHandlerFunc func(w http.ResponseWriter, r *http.Request, recovery interface{})

func defaultRecoveryHandler(w http.ResponseWriter, r *http.Request, recovery interface{}) {
	w.WriteHeader(http.StatusInternalServerError)

	log.Println(recovery)
	debug.PrintStack()
}

type recoveryOptions struct {
	handlerFunc RecoveryHandlerFunc
}

// RecoveryOption An Option that config RecoverMiddleware
type RecoveryOption interface {
	apply(opts *recoveryOptions)
}

type recoveryOptionFunc func(opts *recoveryOptions)

func (optFunc recoveryOptionFunc) apply(opts *recoveryOptions) {
	optFunc(opts)
}

// WithRecoveryHandlerOption Custom recovery handler
func WithRecoveryHandlerOption(handlerFunc RecoveryHandlerFunc) RecoveryOption {
	return recoveryOptionFunc(func(opts *recoveryOptions) {
		opts.handlerFunc = handlerFunc
	})
}

type recoveryMiddleware struct {
	next http.Handler
	opts recoveryOptions
}

func (middleware *recoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		recovery := recover()
		if recovery != nil {
			middleware.handleRecovery(w, r, recovery)
		}
	}()

	middleware.next.ServeHTTP(w, r)
}

func (middleware *recoveryMiddleware) handleRecovery(w http.ResponseWriter, r *http.Request, recovery interface{}) {
	if middleware.opts.handlerFunc != nil {
		middleware.opts.handlerFunc(w, r, recovery)
	} else {
		defaultRecoveryHandler(w, r, recovery)
	}
}

// RecoverMiddleware Create a HTTP middleware that recover panic
func RecoverMiddleware(opts ...RecoveryOption) Middleware {
	return func(next http.Handler) http.Handler {
		middleware := recoveryMiddleware{
			next: next,
		}
		for _, opt := range opts {
			opt.apply(&middleware.opts)
		}
		return &middleware
	}
}
