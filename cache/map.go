package cache

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// MapCache - simplest cache implementation based on map and mutex
type MapCache struct {
	store      map[uint64]*Item
	lock       sync.RWMutex
	expiration time.Duration
}

// Get - deliver cached item if it exists
func (c *MapCache) Get(key uint64) (*Item, bool) {
	c.lock.RLock()
	item, ok := c.store[key]
	c.lock.RUnlock()
	return item, ok
}

// Set - set cached item
func (c *MapCache) Set(key uint64, item *Item) {
	c.lock.Lock()
	c.store[key] = item
	c.lock.Unlock()
}

// CanBeCached - decide if we can cache that request
func (c *MapCache) CanBeCached(key uint64) bool {
	return true
}

// GatherData - collect response header data for future decisions
func (c *MapCache) GatherData(key uint64, header http.Header) {}

// GarbageCollect - each 30 seconds clear stale items,
// should be started as goroutine
func (c *MapCache) GarbageCollect() {
	for {
		select {
		case <-time.After(30 * time.Second):
			c.lock.Lock()
			newMap := make(map[uint64]*Item)
			for key, item := range c.store {
				if !item.isExpired(c.expiration) {
					newMap[key] = item
				}
			}
			c.store = newMap
			c.lock.Unlock()
			log.Println("[cache] running garbage collection")
		}
	}
}
