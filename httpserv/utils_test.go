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
