package middlewares_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wencan/middlewares"
)

func TestRecoveryMiddleware(t *testing.T) {
	panicError := "this is a test"
	panicHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		panic(panicError)
	}

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()
	middleware := middlewares.RecoverMiddleware()
	handler := middleware(http.HandlerFunc(panicHandlerFunc))
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestRecoveryMiddlewareWithHandler(t *testing.T) {
	panicError := "this is a test"
	panicHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		panic(panicError)
	}

	var buffer bytes.Buffer
	handlerFunc := func(w http.ResponseWriter, r *http.Request, recovery interface{}) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(&buffer, recovery)
	}

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()
	middleware := middlewares.RecoverMiddleware(middlewares.WithRecoveryHandlerOption(handlerFunc))
	handler := middleware(http.HandlerFunc(panicHandlerFunc))
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	assert.Equal(t, panicError, buffer.String())
}
