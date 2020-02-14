package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// RoundRobinBucket - round-robin representation of targets bucket
type RoundRobinBucket struct {
	targets []Target     // targets storage
	last    uint64       // last used target index
	lock    sync.RWMutex // lock for target slice
}

// Addtarget - add target to storage
func (buck *RoundRobinBucket) AddTarget(trg Target) error {
	if trg == nil {
		return ErrInvalidTarget
	}
	trg.ReverseProxy().SetErrHandler(buck.getErrHandler(trg))
	status := trg.Ping()
	if !status {
		log.Println("[loadbalancer]", trg.Address(), "unreachable")
	}
	trg.SetAvailable(status)
	buck.lock.Lock()
	buck.targets = append(buck.targets, trg)
	buck.lock.Unlock()
	return nil
}

// Size - amount of stored targets
func (buck *RoundRobinBucket) Size() int {
	buck.lock.RLock()
	defer buck.lock.RUnlock()
	return len(buck.targets)
}

// ServeHTTP - serve incoming request with target's proxy
func (buck *RoundRobinBucket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trg, err := buck.getNextTarget()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[loadbalancer]", err)
		return
	}
	log.Println("[loadbalancer] choose", trg.Address())
	proxy := trg.ReverseProxy()
	// log.Println("[proxy] to", trg.Address())
	proxy.ServeHTTP(w, r)
	return
}

// getNextTarget - round-robin algorithm for chosing next target
// Check if current target is available to serve request,
// in other case we just get next, while not find good one
func (buck *RoundRobinBucket) getNextTarget() (Target, error) {
	trgAmount := uint64(len(buck.targets))
	if trgAmount == 0 {
		return nil, ErrNoTargetsAvailable
	}
	next := buck.last % trgAmount
	defer atomic.AddUint64(&buck.last, 1)
	full := next + trgAmount
	buck.lock.Lock()
	for pos := next; pos < full; pos++ {
		index := pos % trgAmount
		if buck.targets[index].IsAvailable() {
			if index != next {
				atomic.StoreUint64(&buck.last, pos)
			}
			trg := buck.targets[index]
			buck.lock.Unlock()
			return trg, nil
		}
	}
	buck.lock.Unlock()
	return nil, ErrAllTargetsUnreachable
}

// Healthcheck - passive target's availability checks
func (buck *RoundRobinBucket) Healthcheck() {
	if buck.Size() < 1 {
		log.Printf("[healthcheck] %s \n", ErrNoTargetsAvailable.Error())
	}
	for _, trg := range buck.targets {
		msg := "available"
		status := trg.Ping()
		trg.SetAvailable(status)
		if !status {
			msg = "unreachable"
		}
		log.Printf("[healthcheck] %s (%s)\n", trg.Address(), msg)
	}
}

// RemoveStale - remove stale targets from storage
func (buck *RoundRobinBucket) RemoveStale(timeout time.Duration) {
	if buck.Size() < 1 {
		return
	}
	buck.lock.Lock()
	newtargets := []Target{}
	for _, trg := range buck.targets {
		addr := trg.Address()
		timeDiff := time.Since(time.Unix(trg.LastSeen(), 0))
		if !trg.IsAvailable() && timeDiff > timeout {
			log.Printf("[remove] %s is stale and will be removed\n", addr)
			continue
		}
		newtargets = append(newtargets, trg)
	}
	if len(newtargets) != len(buck.targets) {
		buck.targets = newtargets
	}
	buck.lock.Unlock()
}

// getErrHandler - error handler func for reverse proxy instance
// First, we try MAX_RETRIES time to serve request with current server
// Second, we recurrently call Serve func, to switch server
// Count retries for each server separately
// Count attempts for each request
func (buck *RoundRobinBucket) getErrHandler(trg Target) func(w http.ResponseWriter, r *http.Request, e error) {
	return func(w http.ResponseWriter, r *http.Request, e error) {
		attempts := GetAttemptsFromContext(r)
		if attempts > maxAttempts {
			log.Printf("[attempt] %s (%s) Too much attempts, refusing\n", r.RemoteAddr, r.URL.Path)
			http.Error(w, ErrTargetUnavailable.Error(), http.StatusServiceUnavailable)
			return
		}
		log.Printf("[loadbalancer] %s %s\n", trg.Address(), e.Error())
		retries := GetRetriesFromContext(r)
		proxy := trg.ReverseProxy()
		if retries < maxRetries {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(r.Context(), RetriesKey, retries+1)

				log.Printf("[retry] %s (%s) Retrying server %d\n", r.RemoteAddr, r.URL.Path, attempts)
				proxy.ServeHTTP(w, r.WithContext(ctx))
			}
			return
		}
		trg.SetAvailable(false)
		log.Printf("[attempt] %s (%s) Attempting server %d\n", r.RemoteAddr, r.URL.Path, attempts)
		ctx := context.WithValue(r.Context(), AttemptsKey, attempts+1)
		buck.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RunServices - execute targets pool services
func (buck *RoundRobinBucket) RunServices(staleTimeout int) {
	go func() {
		for {
			select {
			case <-time.After(healthCheckPeriod):
				buck.Healthcheck()
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(removeStalePeriod):
				buck.RemoveStale(time.Minute * time.Duration(staleTimeout))
			}
		}
	}()
}
