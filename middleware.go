package middlewares

import "net/http"

type MiddleWare func(http.Handler) http.Handler

type Chain func(middlewares ...MiddleWare) MiddleWare
