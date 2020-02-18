package cache

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LimitedMapCache - limited cache implementation based on map and mutex
type LimitedMapCache struct {
	maxSize     int
	maxItemSize int

	size     int
	sizeLock sync.RWMutex

	forStore map[uint64]bool
	forLock  sync.RWMutex

	store map[uint64]*Item
	lock  sync.RWMutex

	expiration time.Duration
}

// Get - deliver cached item if it exists
func (c *LimitedMapCache) Get(key uint64) (*Item, bool) {
	c.lock.RLock()
	item, ok := c.store[key]
	c.lock.RUnlock()
	return item, ok
}

// Set - set cached item
func (c *LimitedMapCache) Set(key uint64, item *Item) {
	if len(item.data) > c.maxItemSize {
		log.Println("[cache] item is too huge: expected max", c.maxItemSize, ", got", len(item.data))
		return
	}
	if c.size+len(item.data) > c.maxSize {
		log.Println("[cache] limit exceeded", c.size, len(item.data), c.maxSize)
		return
	}
	if len(item.data) == 0 {
		log.Println("[cache] item has no data, why should we cache it?")
		return
	}
	c.sizeLock.Lock()
	c.size += len(item.data)
	c.sizeLock.Unlock()

	c.lock.Lock()
	c.store[key] = item
	c.lock.Unlock()
	log.Println("[cache] size", c.size)
}

// CanBeCached - decide if we can cache that request
func (c *LimitedMapCache) CanBeCached(key uint64) bool {
	val, ok := c.forStore[key]
	if !ok {
		return false
	}
	return val
}

// GatherData - collect response header data for future decisions
func (c *LimitedMapCache) GatherData(key uint64, header http.Header) {
	contentLength := header.Get("Content-Length")
	itemSize, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Println("[cache] invalid Content-Length", contentLength)
		return
	}
	if itemSize > c.maxItemSize {
		log.Println("[cache] item is too huge: expected max", c.maxItemSize, ", got", itemSize, "for", key)
		return
	}
	if itemSize == 0 {
		log.Println("[cache] there is no content, why should we cache it?")
		return
	}
	controls := strings.Split(header.Get("Cache-Control"), ",")
	for _, item := range controls {
		if item == "private" || item == "no-cache" {
			return
		}
	}
	c.forLock.Lock()
	c.forStore[key] = true
	c.forLock.Unlock()
}

// GarbageCollect - each 30 seconds clear stale items,
// should be started as goroutine
func (c *LimitedMapCache) GarbageCollect() {
	for {
		select {
		case <-time.After(30 * time.Second):
			log.Println("[cache] running garbage collection")
			c.lock.Lock()
			c.sizeLock.Lock()
			c.forLock.Lock()
			log.Println("[cache] current size", c.size)
			newSize := 0
			newMap := make(map[uint64]*Item)
			for key, item := range c.store {
				if !item.isExpired(c.expiration) {
					newMap[key] = item
					newSize += len(item.data)
					continue
				}
				delete(c.forStore, key)
			}
			c.size = newSize
			c.store = newMap
			log.Println("[cache] size after GC", c.size)
			c.forLock.Unlock()
			c.sizeLock.Unlock()
			c.lock.Unlock()
		}
	}
}
