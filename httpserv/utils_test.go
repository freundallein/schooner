package httpserv

import (
	"net/http"
	"testing"
)

func TestGetClientIPReal(t *testing.T) {
	request := &http.Request{Header: http.Header{}}
	expected := "127.0.0.1"
	request.Header.Set("X-Real-Ip", expected)
	observed := GetClientIP(request)
	if observed != expected {
		t.Error("expected", expected, "got", observed)
	}
}
func TestGetClientIPFrowarded(t *testing.T) {
	request := &http.Request{Header: http.Header{}}
	expected := "127.0.0.1"
	request.Header.Set("X-Forwarded-For", expected)
	observed := GetClientIP(request)
	if observed != expected {
		t.Error("expected", expected, "got", observed)
	}
}
func TestGetClientIPRemote(t *testing.T) {
	expected := "127.0.0.1"
	request := &http.Request{Header: http.Header{}, RemoteAddr: expected}
	observed := GetClientIP(request)
	if observed != expected {
		t.Error("expected", expected, "got", observed)
	}
}

func TestConstructKey(t *testing.T) {
	uri := "test"
	key := СonstructKey(uri)
	if key != 18007334074686647077 {
		t.Error("expected 18007334074686647077, got", key)
	}
}

func TestNotCachableFalse(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	observed := NotCachable(req)
	if observed {
		t.Error("expected false, got true")
	}
}
func TestNotCachableTrue(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Error(err)
	}
	observed := NotCachable(req)
	if !observed {
		t.Error("expected true, got false")
	}
}
