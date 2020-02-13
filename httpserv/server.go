package httpserv

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrNoOptions = errors.New("no options provided")
)

// Options - fileserver parameters
type Options struct {
	Port string
}

// String - simple representation
func (opts *Options) String() string {
	return fmt.Sprintf("PORT=%s", opts.Port)
}

// Server - main control struct
type Server struct {
	mux     *http.ServeMux
	options *Options
}

// New - service constructor
func New(options *Options) (*Server, error) {
	if options == nil {
		return nil, ErrNoOptions
	}
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	mux := http.NewServeMux()
	mux.Handle("/schooner/metrics", promhttp.Handler())
	mux.HandleFunc("/schooner/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/", MiddlewareChain(
		dummyHandler,
		AccessLog,
	))
	return &Server{options: options, mux: mux}, nil
}

// Run - start schooner
func (srv *Server) Run() error {
	addr := fmt.Sprintf("0.0.0.0:%s", srv.options.Port)
	serv := &http.Server{
		Handler:        srv.mux,
		Addr:           addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[httpserv] Start listening on %s\n", addr)
	err := serv.ListenAndServe()
	return err
}
