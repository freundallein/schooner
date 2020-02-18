package httpserv

import (
	"hash/fnv"
	"net/http"
)

func GetClientIP(r *http.Request) string {
	addr := r.Header.Get("X-Real-Ip")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = r.RemoteAddr
	}
	return addr
}

// СonstructKey - make hash from uri
func СonstructKey(uri string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(uri))
	return hash.Sum64()
}

// NotCachable - sift non-cachable requests
func NotCachable(r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return false
	}
	return true
}
