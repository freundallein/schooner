package cache

import (
	"net/http"
)

// Cache - basic cache interface
type Cache interface {
	Get(key uint64) (*Item, bool)
	Set(key uint64, item *Item)

	CanBeCached(uint64) bool
	GatherData(uint64, http.Header)

	GarbageCollect()
}
