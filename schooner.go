package main

import (
	"fmt"
	"github.com/freundallein/schooner/httpserv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "02.01.2006 15:04:05"

	machineIDKey     = "MACHINE_ID"
	defaultMachineID = 0

	portKey     = "PORT"
	defaultPort = "8000"

	targetsKey     = "TARGETS"
	defaultTargets = ""

	staleTimeoutKey     = "STALE_TIMEOUT"
	defaultStaleTimeout = 60 // seconds

	useCacheKey     = "USE_CACHE"
	defaultUseCache = 1

	cacheExpireKey     = "CACHE_EXPIRE"
	defaultCacheExpire = 60

	maxCacheSizeKey     = "MAX_CACHE_SIZE"
	defaultMaxCacheSize = 1024 * 1024 // kb

	maxCacheItemSizeKey     = "MAX_CACHE_ITEM_SIZE"
	defaultMaxCacheItemSize = 10 * 1024 // kb
)

type logWriter struct{}

// Write - custom logger formatting
func (writer logWriter) Write(bytes []byte) (int, error) {
	msg := fmt.Sprintf("%s | [schooner] %s", time.Now().UTC().Format(timeFormat), string(bytes))
	return fmt.Print(msg)
}

func getEnv(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return fallback, nil
}

func getIntEnv(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fallback, err
		}
		return int(i), nil
	}
	return fallback, nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	port, err := getEnv(portKey, defaultPort)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	configTargets, err := getEnv(targetsKey, defaultTargets)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	targets := strings.Split(configTargets, ";")
	if len(targets) == 0 {
		log.Fatal("[config] No targets provided")
	}
	staleTimeout, err := getIntEnv(staleTimeoutKey, defaultStaleTimeout)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	machineID, err := getIntEnv(machineIDKey, defaultMachineID)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	useCache, err := getIntEnv(useCacheKey, defaultUseCache)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	cacheExpire, err := getIntEnv(cacheExpireKey, defaultCacheExpire)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	maxCacheSize, err := getIntEnv(maxCacheSizeKey, defaultMaxCacheSize)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	maxCacheItemSize, err := getIntEnv(maxCacheItemSizeKey, defaultMaxCacheItemSize)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}

	options := &httpserv.Options{
		Port:             port,
		Targets:          targets,
		StaleTimeout:     staleTimeout,
		MachineID:        machineID,
		UseCache:         useCache,
		CacheExpire:      cacheExpire,
		MaxCacheSize:     maxCacheSize,
		MaxCacheItemSize: maxCacheItemSize,
	}
	log.Printf("[config] starting with %s\n", options)
	srv, err := httpserv.New(options)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("[httpserv] %s\n", err.Error())
	}
}
