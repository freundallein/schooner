package httpserv

import (
	"log"
	"net/http"
	"strconv"

	"github.com/freundallein/schooner/corridgen"
)

// Middleware - http middleware
type Middleware func(http.Handler) http.Handler

// MiddlewareChain - chain multiple middlewares
func MiddlewareChain(handler http.Handler, midllewares ...Middleware) http.Handler {
	if len(midllewares) < 1 {
		return handler
	}
	wrapped := handler
	for i := 0; i < len(midllewares); i++ {
		wrapped = midllewares[i](wrapped)
	}
	return wrapped
}

// AccessLog - log client requests
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)
		log.Printf("[server] [%s] %s %s \n", clientIP, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// EnrichCorrelationID - add unique correlation id to header
func EnrichCorrelationID(gen *corridgen.Generator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("Correlation-Id", strconv.FormatUint(gen.GetID(), 10))
			next.ServeHTTP(w, r)
		})
	}
}
