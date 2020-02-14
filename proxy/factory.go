package proxy

import (
	"net/http"
	"net/url"
)

// New - reverse proxy factory
func New(strategy ProxyStrategy, addr *url.URL) Proxy {
	switch strategy {
	case DefaultStrategy:
		return &DefaultProxy{
			transport: http.DefaultTransport,
			addr:      addr,
		}
	}
	return nil
}
