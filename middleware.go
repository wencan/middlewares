package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

type Chain func(middlewares ...Middleware) Middleware
