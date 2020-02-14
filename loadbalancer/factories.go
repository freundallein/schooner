package loadbalancer

import (
	"net/url"
	"time"

	"github.com/freundallein/schooner/proxy"
)

// NewTarget - backend target factory
func NewTarget(URL string) (Target, error) {
	addr, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	reverseProxy := proxy.New(proxy.DefaultStrategy, addr)
	return &DefaultTarget{
		address:      addr,
		isAvailable:  true,
		reverseProxy: reverseProxy,
		lastSeen:     time.Now().Unix(),
	}, nil
}

// New - backends pool factory, can use different balancing algorithms in future
func New(algo LoadBalanceAlgorithm) (TargetBucket, error) {
	var bckt TargetBucket
	switch algo {
	case RoundRobin:
		bckt = &RoundRobinBucket{
			targets: []Target{},
		}
	}
	if bckt == nil {
		return nil, ErrInvalidAlgorithm
	}
	return bckt, nil
}
