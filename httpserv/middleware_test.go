package httpserv

import (
	"net/http"
	"testing"
)

type MockHandler struct{}

func (mh *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestMiddlewareChainEmpty(t *testing.T) {
	handler := http.Handler(&MockHandler{})
	chain := MiddlewareChain(handler)
	if chain != handler {
		t.Error("Expected to receive same handler")
	}
}
