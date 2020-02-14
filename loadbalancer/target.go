package loadbalancer

import (
	"github.com/freundallein/schooner/proxy"
	"net"
	"net/url"
	"sync"
	"time"
)

// DefaultTarget - default backend target implementation
type DefaultTarget struct {
	address      *url.URL     // target address
	isAvailable  bool         // current status
	lock         sync.RWMutex // lock for isAvailable attribute
	reverseProxy proxy.Proxy  // reverse proxy for request forwarding
	lastSeen     int64        // unixtime for last time, when server was available
}

// IsAvailable - getter for server's availability
func (dt *DefaultTarget) IsAvailable() bool {
	dt.lock.RLock()
	status := dt.isAvailable
	dt.lock.RUnlock()
	return status
}

// SetAvailable - setter for server's availability
func (dt *DefaultTarget) SetAvailable(status bool) {
	dt.lock.Lock()
	dt.isAvailable = status
	if status {
		dt.lastSeen = time.Now().Unix()
	}
	dt.lock.Unlock()
}

// Address - getter for server address
func (dt *DefaultTarget) Address() *url.URL {
	return dt.address
}

// ReverseProxy - getter for reversing proxy
func (dt *DefaultTarget) ReverseProxy() proxy.Proxy {
	return dt.reverseProxy
}

// LastSeen - getter for lastSeen time field
func (dt *DefaultTarget) LastSeen() int64 {
	return dt.lastSeen
}

func (dt *DefaultTarget) Ping() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", dt.address.Host, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
