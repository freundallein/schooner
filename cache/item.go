package cache

import (
	"net/http"
	"time"
)

// Item - basic cache item representation
type Item struct {
	statusCode int
	header     http.Header
	data       []byte
	lastSeen   time.Time
}

// NewItem - constructor
func NewItem(statusCode int, header http.Header, data []byte) *Item {
	return &Item{
		header:     header,
		statusCode: statusCode,
		data:       data,
		lastSeen:   time.Now().UTC(),
	}
}

// StatusCode - status code getter
func (i *Item) StatusCode() int {
	return i.statusCode
}

// Header - header data getter
func (i *Item) Header() http.Header {
	return i.header
}

// Data - item's data getter
func (i *Item) Data() []byte {
	return i.data
}

//isExpired - check if current item is stale
func (i *Item) isExpired(expiration time.Duration) bool {
	return time.Since(i.lastSeen) > expiration
}
