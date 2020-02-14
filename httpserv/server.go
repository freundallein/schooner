package httpserv

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/freundallein/schooner/loadbalancer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrNoOptions = errors.New("no options provided")
)

// Options - fileserver parameters
type Options struct {
	Port         string
	Targets      []string
	StaleTimeout int
	MachineID    int
	UseCache     int
	CacheExpire  int
}

// String - simple representation
func (opts *Options) String() string {
	return fmt.Sprintf(
		"PORT=%s, TARGETS=%s, STALE_TIMEOUT=%d, MACHINE_ID=%d, USE_CACHE=%d, CACHE_EXPIRE=%d",
		opts.Port, opts.Targets, opts.StaleTimeout, opts.MachineID, opts.UseCache, opts.CacheExpire,
	)
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
	bucket, err := loadbalancer.New(loadbalancer.RoundRobin)
	if err != nil {
		return nil, err
	}
	for _, target := range options.Targets {
		trg, err := loadbalancer.NewTarget(target)
		if err != nil {
			return nil, err
		}
		bucket.AddTarget(trg)
		log.Printf("[config] target %s added\n", target)
	}

	mux := http.NewServeMux()
	mux.Handle("/schooner/metrics", promhttp.Handler())
	mux.HandleFunc("/schooner/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/", MiddlewareChain(
		bucket,
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
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[httpserv] Start listening on %s\n", addr)
	err := serv.ListenAndServe()
	return err
}
