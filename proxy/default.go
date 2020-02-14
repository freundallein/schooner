package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// DefaultProxy - basic proxy via http.DefaultTransport
type DefaultProxy struct {
	addr       *url.URL
	transport  http.RoundTripper
	ErrHandler func(http.ResponseWriter, *http.Request, error)
}

// ServeHTTP - handle http request
func (prx *DefaultProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	patchTargetAddr(prx.addr, req)
	response, err := prx.transport.RoundTrip(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		log.Println("[proxy]", err)
		prx.handleError(w, req, err)
		return
	}
	transferHeaders(w.Header(), response.Header)
	w.WriteHeader(response.StatusCode)
	err = transferBody(w, response.Body)
	if err != nil {
		log.Println("[proxy]", err)
		prx.handleError(w, req, err)
		return
	}
}

// SetErrHandler - set
func (prx *DefaultProxy) SetErrHandler(f func(w http.ResponseWriter, req *http.Request, err error)) {
	prx.ErrHandler = f
}

func (prx *DefaultProxy) handleError(w http.ResponseWriter, req *http.Request, err error) {
	if prx.ErrHandler == nil {
		return
	}
	prx.ErrHandler(w, req, err)
}

func patchTargetAddr(addr *url.URL, req *http.Request) {
	req.URL.Scheme = addr.Scheme
	req.URL.Host = addr.Host
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "")
	}
}

func transferHeaders(to, from http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func transferBody(w http.ResponseWriter, response io.Reader) error {
	_, err := io.Copy(w, response)
	return err
}
