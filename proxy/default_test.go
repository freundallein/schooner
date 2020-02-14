package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HELLO"))
	}))
	defer ts.Close()
	target, _ := url.Parse(ts.URL)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888", nil)
	if err != nil {
		t.Error(err)
	}
	prx := &DefaultProxy{addr: target, transport: http.DefaultTransport}
	rec := httptest.NewRecorder()
	prx.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	resp := rec.Result()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "HELLO" {
		t.Error("Expected HELLO got", string(body))
	}
}

func TestSetErrHandler(t *testing.T) {
	prx := &DefaultProxy{}
	f := func(w http.ResponseWriter, req *http.Request, err error) {}
	expected := reflect.ValueOf(f)
	prx.SetErrHandler(f)
	if reflect.ValueOf(prx.ErrHandler) != expected {
		t.Error("expected", expected, "got", reflect.ValueOf(prx.ErrHandler))
	}
}

func TesthandleError(t *testing.T) {
	prx := &DefaultProxy{}
	flag := false
	f := func(w http.ResponseWriter, req *http.Request, err error) {
		flag = true
	}
	request := &http.Request{}
	writer := httptest.NewRecorder()
	prx.handleError(writer, request, nil)
	if flag {
		t.Error("Expected not calling nil error handler")
	}
	prx.SetErrHandler(f)
	prx.handleError(writer, request, nil)
	if !flag {
		t.Error("Expected calling error handler")
	}
}
func TestPatchTargetAddr(t *testing.T) {
	original, _ := url.Parse("http://localhost:8000/")
	target, _ := url.Parse("http://localhost:8001/")
	request := &http.Request{URL: original, Header: make(http.Header)}
	patchTargetAddr(target, request)
	if request.URL.Scheme != target.Scheme {
		t.Error("expected", target.Scheme, "got", request.URL.Scheme)
	}
	if request.URL.Host != target.Host {
		t.Error("expected", target.Host, "got", request.URL.Host)
	}
	if val, ok := request.Header["User-Agent"]; !ok || val[0] != "" {
		t.Error("expected User-Agent header")
	}

}

func TestTransferHeaders(t *testing.T) {
	from := make(http.Header)
	from.Add("test", "1")
	from.Add("test", "2")
	from.Add("test", "3")
	from.Add("source", "123456")
	to := make(http.Header)
	transferHeaders(to, from)
	for key, val := range from {
		observed := to[key]
		for idx, item := range val {
			if item != observed[idx] {
				t.Error("expected", item, "got", observed[idx])
			}
		}
	}
}

func TestTransferBody(t *testing.T) {
	data := `{"ap_id":"123","probe_requests":[{"mac":"1","timestamp":"2","bssid":"3","ssid":"4"}]}`
	req, err := http.NewRequest("GET", "http://127.0.0.1:8888", bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Error(err)
	}
	writer := httptest.NewRecorder()
	err = transferBody(writer, req.Body)
	if err != nil {
		t.Error(err)
	}
	resp := writer.Result()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != data {
		t.Error("Expected", data, "got", string(body))
	}

}
