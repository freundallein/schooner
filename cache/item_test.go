package cache

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	observed := NewItem(200, make(http.Header), []byte{})
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&Item{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestStatusCode(t *testing.T) {
	observed := NewItem(200, make(http.Header), []byte{})
	if observed.StatusCode() != 200 {
		t.Error("Expected 200 got", observed.StatusCode())
	}
}

func TestHeader(t *testing.T) {
	hdr := make(http.Header)
	hdr.Add("Key", "Value")
	observed := NewItem(1, make(http.Header), []byte{})
	if observed.Header() == nil {
		t.Error("Expected header got nil")
	}
	for key, value := range observed.Header() {
		if value[0] != hdr.Get(key) {
			t.Error("Expected", hdr.Get(key), "got", value, "for", key)
		}
	}
}

func TestData(t *testing.T) {
	expected := []byte{1, 2, 3}
	observed := NewItem(200, make(http.Header), expected).Data()
	if observed == nil {
		t.Error("Expected data got nil")
	}
	for idx, item := range observed {
		if item != expected[idx] {
			t.Error("Expected", expected[idx], "got", item)
		}
	}
}

func TestIsExpired(t *testing.T) {
	item := NewItem(200, make(http.Header), []byte{})
	time.Sleep(1 * time.Millisecond)
	observed := item.isExpired(2 * time.Millisecond)
	if observed == true {
		t.Error("Expect false")
	}
	time.Sleep(2 * time.Millisecond)
	observed = item.isExpired(2 * time.Millisecond)
	if observed == false {
		t.Error("Expect true")
	}
}

// //isExpired - check if current item is stale
// func (i *Item) isExpired(expiration time.Duration) bool {
// 	return time.Since(i.lastSeen) > expiration
// }
