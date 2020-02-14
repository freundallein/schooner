package httpserv

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewWithError(t *testing.T) {
	observed, err := New(nil)
	if err == nil {
		t.Error("expected", ErrNoOptions.Error(), "got", err.Error())
	}
	if observed != nil {
		t.Error("expected", nil, "got", observed)
	}
}
func TestNew(t *testing.T) {
	observed, err := New(&Options{})
	if err != nil {
		t.Error(err.Error())
	}
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&Server{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestHealthz(t *testing.T) {
	srv, err := New(&Options{})
	if err != nil {
		t.Error(err.Error())
	}
	req, err := http.NewRequest("GET", "/schooner/healthz", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()
	srv.mux.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}
	resp := rec.Result()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "OK" {
		t.Error("Expected OK got", string(body))
	}
}
func TestMetrics(t *testing.T) {
	srv, err := New(&Options{})
	if err != nil {
		t.Error(err.Error())
	}
	req, err := http.NewRequest("GET", "/schooner/metrics", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()
	srv.mux.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}
}

func TestIndexNoTargets(t *testing.T) {
	srv, err := New(&Options{})
	if err != nil {
		t.Error(err.Error())
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()
	srv.mux.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusInternalServerError)
	}
}
