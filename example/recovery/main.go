package main

import (
	"log"
	"net/http"

	"github.com/wencan/middlewares"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Println(err)
		return
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		panic("this is a test")
	}

	middleware := middlewares.RecoverMiddleware(middlewares.WithRecoveryHandlerOption(func(w http.ResponseWriter, r *http.Request, recovery interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("recover a panic", zap.Any("panic", recovery))
	}))

	handler := middleware(http.HandlerFunc(handlerFunc))
	err = http.ListenAndServe(":8080", handler)
	if err != nil && err != http.ErrServerClosed {
		logger.Error("server closed", zap.Error(err))
	}
}
