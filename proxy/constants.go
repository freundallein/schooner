package proxy

type ProxyStrategy string

const (
	// DefaultStrategy - simple proxy with http.Transport
	DefaultStrategy ProxyStrategy = "default"
)
