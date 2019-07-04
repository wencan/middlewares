package main

import (
	"log"
	"net/http"

	"github.com/wencan/middlewares"
	middleware_zap "github.com/wencan/middlewares/zap"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Println(err)
		return
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	}

	logger := middleware_zap.NewLogger(zapLogger)
	middleware := middlewares.LoggingMiddleware(logger)

	handler := middleware(http.HandlerFunc(handlerFunc))
	err = http.ListenAndServe(":8080", handler)
	if err != nil && err != http.ErrServerClosed {
		zapLogger.Error("server closed", zap.Error(err))
	}
}
