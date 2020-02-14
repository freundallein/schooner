package loadbalancer

import (
	"errors"
	"time"
)

type LoadBalanceAlgorithm string

const (
	// Available loadbalancing algorithms
	RoundRobin LoadBalanceAlgorithm = "round-robin"

	// ErrorHandler context keys
	AttemptsKey = "attempts"
	RetriesKey  = "retries"
	maxRetries  = 3
	maxAttempts = 3

	// Service periods
	healthCheckPeriod = 5 * time.Second
	removeStalePeriod = 60 * time.Second
)

var (
	// Custom errors
	ErrInvalidAlgorithm      = errors.New("invalid balancing algorithm chosen.")
	ErrInvalidTarget         = errors.New("expected Target, got nil")
	ErrNoTargetsAvailable    = errors.New("no targets available")
	ErrAllTargetsUnreachable = errors.New("all targets unreachable")
	ErrTargetUnavailable     = errors.New("target not available")
)
