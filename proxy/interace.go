package proxy

import (
	"net/http"
)

// Proxy - basic reverse proxy interface
type Proxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	SetErrHandler(func(http.ResponseWriter, *http.Request, error))
}
