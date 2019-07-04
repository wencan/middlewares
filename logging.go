package middlewares

/*
 * http logging handler
 *
 * wencan
 * 2019-06-27
 */

import (
	"net/http"
	"time"
)

// LoggingLogger Interface of LoggingMiddleware logger
type LoggingLogger interface {
	Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error
}

// NopLoggingLogger Implement of LoggingLogger. It do nothing.
type NopLoggingLogger struct {
}

func (logger *NopLoggingLogger) Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error {
	return nil
}

type loggingOptions struct {
}

// LoggingOption An Option that config LoggingMiddleware
type LoggingOption interface {
	apply(opts *loggingOptions)
}

type loggingOptionFunc func(opts *loggingOptions)

func (optFunc loggingOptionFunc) apply(opts *loggingOptions) {
	optFunc(opts)
}

type loggingMiddleware struct {
	logger LoggingLogger
	opts   loggingOptions

	next http.Handler
}

func (middleware *loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &response{w: w}
	timestamp := time.Now()

	middleware.next.ServeHTTP(resp, r)

	middleware.logger.Write(r, resp.Status(), resp.BodyBytesSent(), timestamp)
}

// LoggingMiddleware Create a HTTP middleware that log request and response
func LoggingMiddleware(logger LoggingLogger, opts ...LoggingOption) Middleware {
	return func(next http.Handler) http.Handler {
		middleware := &loggingMiddleware{
			logger: logger,
			next:   next,
		}
		for _, opt := range opts {
			opt.apply(&middleware.opts)
		}
		return middleware
	}
}
