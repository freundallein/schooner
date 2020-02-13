package httpserv

import (
	"net/http"
	"testing"
)

// import (
// 	"log"
// 	"net/http"
// 	"strings"
// )

// func MockHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// }

type MockHandler struct{}

func (mh *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func TestMiddlewareChainEmpty(t *testing.T) {
	handler := http.Handler(&MockHandler{})
	chain := MiddlewareChain(handler)
	if chain != handler {
		t.Error("Expected to receive same handler")
	}
}
