package cache

import (
	"time"
)

// CacheStrategy
type CacheStrategy string

const (
	// InfiniteMapStrategy - use simple map without any restrictions
	// Is bad idea for all cases
	InfiniteMapStrategy CacheStrategy = "infinite-map"
	// LimitedMapStrategy - us map with limited size and limited item size
	// No cache eviction except GC
	LimitedMapStrategy CacheStrategy = "limited-map"
	// TODO: LimitedLRUStrategy  CacheStrategy = "lru"
)

// New - cache factory
func New(strategy CacheStrategy, expiration time.Duration, maxSize, maxItemSize int) Cache {
	switch strategy {
	case InfiniteMapStrategy:
		return &MapCache{
			store:      map[uint64]*Item{},
			expiration: expiration,
		}
	case LimitedMapStrategy:
		return &LimitedMapCache{
			store:       map[uint64]*Item{},
			forStore:    map[uint64]bool{},
			expiration:  expiration,
			maxSize:     maxSize * 1024,     // kB to bytes
			maxItemSize: maxItemSize * 1024, // kB to bytes
		}
	}
	return nil
}
