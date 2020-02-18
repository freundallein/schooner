package httpserv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/freundallein/schooner/cache"
	"github.com/freundallein/schooner/corridgen"
)

type MockHandler struct {
	data       []byte
	statusCode int
	header     http.Header
}

func (mh *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, values := range mh.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(mh.data)))
	w.WriteHeader(mh.statusCode)
	w.Write(mh.data)
}

func TestMiddlewareChainEmpty(t *testing.T) {
	handler := http.Handler(&MockHandler{})
	chain := MiddlewareChain(handler)
	if chain != handler {
		t.Error("Expected to receive same handler")
	}
}

func TestEnrichCorrelationID(t *testing.T) {
	gen := corridgen.New(uint8(0))
	middleware := EnrichCorrelationID(gen)
	handler := middleware(&MockHandler{})
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if req.Header.Get("Correlation-Id") == "" {
		t.Error("expected Correlation-Id header")
	}

}

func TestCache(t *testing.T) {
	store := cache.New(cache.LimitedMapStrategy, 1*time.Second, 10, 1)
	middleware := Cache(store)
	expectedData := []byte{byte(1), byte(2), byte(3)}
	handler := middleware(&MockHandler{
		data:       expectedData,
		statusCode: 200,
		header:     make(http.Header),
	})
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888", nil)
	if err != nil {
		t.Error(err)
	}
	// First request (gather data)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	// 14695981039346656037 - key for "/"
	if !store.CanBeCached(14695981039346656037) {
		t.Error("Expected true, got false")
	}
	cached, ok := store.Get(14695981039346656037)
	if ok || cached != nil {
		t.Error("Expected empty cache, got", cached)
	}
	// Second request (cache data)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	cached, ok = store.Get(14695981039346656037)
	if !ok || cached == nil {
		t.Error("Expected item in cache, got nil")
	}
	// Third request (respond with cached data)
	// Create new handler to make different response
	newHandler := middleware(&MockHandler{
		data:       []byte{},
		statusCode: 400,
		header:     make(http.Header),
	})
	rec = httptest.NewRecorder()
	newHandler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Error("Expected 200, got", rec.Code)
	}
	if len(rec.Body.Bytes()) != 3 {
		t.Error("Expected 3, got", rec.Body.Bytes())
	}

}
