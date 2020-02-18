package httpserv

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/freundallein/schooner/cache"
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

// Cache - cache requests
func Cache(store cache.Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if NotCachable(r) {
				next.ServeHTTP(w, r)
				return
			}
			uri := r.URL.Path
			cacheKey := СonstructKey(uri)
			if !store.CanBeCached(cacheKey) {
				next.ServeHTTP(w, r)
				store.GatherData(cacheKey, w.Header())
				return
			}
			if response, ok := store.Get(cacheKey); ok {
				for key, values := range response.Header() {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.WriteHeader(response.StatusCode())
				w.Write(response.Data())
				log.Printf("[cache] %s (hit) %d\n", r.URL.Path, response.StatusCode())
				return
			}
			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)
			if recorder.Code < 400 {
				response := cache.NewItem(recorder.Code, recorder.Header(), recorder.Body.Bytes())
				store.Set(cacheKey, response)
			}
			for key, value := range recorder.Header() {
				w.Header().Set(key, strings.Join(value, ","))
			}
			w.WriteHeader(recorder.Code)
			w.Write(recorder.Body.Bytes())
		})
	}
}
