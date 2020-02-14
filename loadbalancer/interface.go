package loadbalancer

import (
	"github.com/freundallein/schooner/proxy"
	"net/http"
	"net/url"
	// "time"
)

// Target - common backend server interface
type Target interface {
	Address() *url.URL
	ReverseProxy() proxy.Proxy

	IsAvailable() bool
	SetAvailable(bool)

	LastSeen() int64

	Ping() bool
}

// TargetBucket - common servers pool interface
type TargetBucket interface {
	AddTarget(Target) error
	ServeHTTP(http.ResponseWriter, *http.Request)
	RunServices(int)
}
