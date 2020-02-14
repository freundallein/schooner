package loadbalancer

import (
	"github.com/freundallein/schooner/proxy"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestIsAvailable(t *testing.T) {
	srv, _ := NewTarget("http://testhost:8000")
	if !srv.IsAvailable() {
		t.Error("Expected true, got false")
	}
}

func TestSetAvailable(t *testing.T) {
	srv, _ := NewTarget("http://testhost:8000")
	srv.SetAvailable(false)
	if srv.IsAvailable() {
		t.Error("Expected false, got true")
	}
	srv.SetAvailable(true)
	if !srv.IsAvailable() {
		t.Error("Expected true, got false")
	}
}

func TestAddress(t *testing.T) {
	srv, _ := NewTarget("http://testhost:8000")
	url := srv.Address()
	if url.Host != "testhost:8000" {
		t.Error("Expected", "testhost:8000", "got", url.Host)
	}
}

func TestReverseProxy(t *testing.T) {
	trg, _ := NewTarget("http://testhost:8000")
	observedType := reflect.TypeOf(trg.ReverseProxy())
	addr, _ := url.Parse("http://localhost:8000/")
	prx := proxy.New(proxy.DefaultStrategy, addr)
	expectedType := reflect.TypeOf(prx)
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestLastSeen(t *testing.T) {
	srv, _ := NewTarget("http://testhost:8000")
	if srv.LastSeen() != time.Now().Unix() {
		t.Error("Expected", time.Now().Unix(), "got", srv.LastSeen())
	}
}
