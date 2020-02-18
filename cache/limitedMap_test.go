package cache

import (
	"net/http"
	"testing"
	"time"
)

// // LimitedMapCache - limited cache implementation based on map and mutex
// type LimitedMapCache struct {
// 	maxSize     int
// 	maxItemSize int

// 	size     int
// 	sizeLock sync.RWMutex

// 	forStore map[uint64]bool
// 	forLock  sync.RWMutex

// 	store map[uint64]*Item
// 	lock  sync.RWMutex

// 	expiration time.Duration
// }

func TestSet(t *testing.T) {
	limitedMap := New(LimitedMapStrategy, 1*time.Second, 10, 1)
	item := NewItem(200, make(http.Header), []byte{1, 2, 3})
	limitedMap.Set(1, item)
	observed, ok := limitedMap.Get(1)
	if !ok {
		t.Error("Expected true, got false")
	}
	if observed != item {
		t.Error("Expected", item, "got", observed)
	}
}
func TestSetHuge(t *testing.T) {
	limitedMap := New(LimitedMapStrategy, 1*time.Second, 10, 1)
	data := []byte{}
	for i := 0; i < 1025; i++ {
		data = append(data, byte(i))
	}
	item := NewItem(200, make(http.Header), data)
	limitedMap.Set(1, item)
	_, ok := limitedMap.Get(1)
	if ok {
		t.Error("Expected false, got true")
	}
}
func TestSetFull(t *testing.T) {
	limitedMap := New(LimitedMapStrategy, 1*time.Second, 1, 10)
	data := []byte{}
	for i := 0; i < 10240; i++ {
		data = append(data, byte(i))
	}
	item := NewItem(200, make(http.Header), data)
	limitedMap.Set(1, item)
	_, ok := limitedMap.Get(1)
	if ok {
		t.Error("Expected false, got true")
	}
}

func TestGetInvalidKey(t *testing.T) {
	limitedMap := New(LimitedMapStrategy, 1*time.Second, 10, 1)
	observed, ok := limitedMap.Get(1)
	if ok {
		t.Error("Expected false, got true")
	}
	if observed != nil {
		t.Error("Expected nil, got", observed)
	}
}

func TestCanBeCached(t *testing.T) {
	limitedMap := &LimitedMapCache{
		store:       map[uint64]*Item{},
		forStore:    map[uint64]bool{},
		expiration:  1 * time.Second,
		maxSize:     10 * 1024, // kB to bytes
		maxItemSize: 1 * 1024,  // kB to bytes
	}
	limitedMap.forStore[1] = true
	observed := limitedMap.CanBeCached(1)
	if !observed {
		t.Error("Expected true, got false")
	}
	observed = limitedMap.CanBeCached(2)
	if observed {
		t.Error("Expected false, got true")
	}
}

func TestGatherData(t *testing.T) {
	limitedMap := &LimitedMapCache{
		store:       map[uint64]*Item{},
		forStore:    map[uint64]bool{},
		expiration:  1 * time.Second,
		maxSize:     10 * 1024, // kB to bytes
		maxItemSize: 1 * 1024,  // kB to bytes
	}
	header := make(http.Header)
	header.Add("Content-Length", "1024")
	limitedMap.GatherData(1, header)
	if val, ok := limitedMap.forStore[1]; !ok || !val {
		t.Error("Expected true, got false")
	}
}
func TestGatherDataHuge(t *testing.T) {
	limitedMap := &LimitedMapCache{
		store:       map[uint64]*Item{},
		forStore:    map[uint64]bool{},
		expiration:  1 * time.Second,
		maxSize:     10 * 1024, // kB to bytes
		maxItemSize: 1 * 1024,  // kB to bytes
	}
	header := make(http.Header)
	header.Add("Content-Length", "1025")
	limitedMap.GatherData(1, header)
	if val, ok := limitedMap.forStore[1]; ok || val {
		t.Error("Expected false, got true")
	}
}
func TestGatherDataCacheControlPrivate(t *testing.T) {
	limitedMap := &LimitedMapCache{
		store:       map[uint64]*Item{},
		forStore:    map[uint64]bool{},
		expiration:  1 * time.Second,
		maxSize:     10 * 1024, // kB to bytes
		maxItemSize: 1 * 1024,  // kB to bytes
	}
	header := make(http.Header)
	header.Add("Content-Length", "1024")
	header.Add("Cache-Control", "private")
	limitedMap.GatherData(1, header)
	if val, ok := limitedMap.forStore[1]; ok || val {
		t.Error("Expected false, got true")
	}
}
func TestGatherDataCacheControlNoCache(t *testing.T) {
	limitedMap := &LimitedMapCache{
		store:       map[uint64]*Item{},
		forStore:    map[uint64]bool{},
		expiration:  1 * time.Second,
		maxSize:     10 * 1024, // kB to bytes
		maxItemSize: 1 * 1024,  // kB to bytes
	}
	header := make(http.Header)
	header.Add("Content-Length", "1024")
	header.Add("Cache-Control", "no-cache")
	limitedMap.GatherData(1, header)
	if val, ok := limitedMap.forStore[1]; ok || val {
		t.Error("Expected false, got true")
	}
}
