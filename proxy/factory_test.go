package proxy

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewDefaultStrategy(t *testing.T) {
	addr, _ := url.Parse("http://localhost:8000/")
	observed := New(DefaultStrategy, addr)
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&DefaultProxy{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
func TestNewWithInvalidStrategy(t *testing.T) {
	addr, _ := url.Parse("http://localhost:8000/")
	const invalidStrategy ProxyStrategy = "whoami"
	observed := New(invalidStrategy, addr)
	if observed != nil {
		t.Error("Expected nil got", observed)
	}
}
