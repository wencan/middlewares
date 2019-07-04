package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/wencan/middlewares"
	middleware_logrus "github.com/wencan/middlewares/logging/logrus"
)

func main() {
	logrusLogger := logrus.New()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	}

	logger := middleware_logrus.NewLogger(logrusLogger)
	middleware := middlewares.LoggingMiddleware(logger)

	handler := middleware(http.HandlerFunc(handlerFunc))
	err := http.ListenAndServe(":8080", handler)
	if err != nil && err != http.ErrServerClosed {
		logrusLogger.WithError(err).Error("server closed")
	}
}
