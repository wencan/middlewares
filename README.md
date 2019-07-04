# middlewares
HTTP middlewares

## LoggingMiddleware

Usage:

```go
import (
    "net/http"
    
	"github.com/wencan/middlewares"
	middleware_zap "github.com/wencan/middlewares/zap"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewDevelopment()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	}

	logger := middleware_zap.NewLogger(zapLogger)
	middleware := middlewares.LoggingMiddleware(logger)

	handler := middleware(http.HandlerFunc(handlerFunc))
	http.ListenAndServe(":8080", handler)
}
```

## RecoveryMiddleware

Usage:

```go
import (
	"net/http"

	"github.com/wencan/middlewares"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		panic("this is a test")
	}

	middleware := middlewares.RecoverMiddleware(middlewares.WithRecoveryHandlerOption(func(w http.ResponseWriter, r *http.Request, recovery interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("recover a panic", zap.Any("panic", recovery))
	}))

	handler := middleware(http.HandlerFunc(handlerFunc))
	http.ListenAndServe(":8080", handler)
}
```

## Chain

Usage:

```go
import (
	"net/http"

	"github.com/wencan/middlewares"
	middleware_zap "github.com/wencan/middlewares/zap"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/panic" {
			panic("BOOM")
		}
		w.Write([]byte("hello, world"))
	}

	middleware := middlewares.Chain(
		middlewares.LoggingMiddleware(middleware_zap.NewLogger(logger)), // logging
		middlewares.RecoverMiddleware(middlewares.WithRecoveryHandlerOption(func(w http.ResponseWriter, r *http.Request, recovery interface{}) {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("recover a panic", zap.Any("panic", recovery))
		})), // recover
	)

	handler := middleware(http.HandlerFunc(handlerFunc))
	http.ListenAndServe(":8080", handler)
}
```