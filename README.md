# schooner
[![Build Status](https://travis-ci.org/freundallein/schooner.svg?branch=master)](https://travis-ci.org/freundallein/schooner)
[![Go Report Card](https://goreportcard.com/badge/github.com/freundallein/schooner)](https://goreportcard.com/report/github.com/freundallein/schooner)

Simple proxy server.  

## Features
* Proxy requests to targets  
* Round-robin load balancing  
* Cache requests  
* Generates unique time-based correlation id for each request (set `Correlation-Id` header)

## Configuration
Application supports configuration via environment variables:
```
PORT=8000 (default 8000)
STALE_TIMEOUT=60 (default 60 - minutes) - used for deleting unavailable targets
MACHINE_ID=0 - used for correlation id generator
USE_CACHE=1  (1 - use / 0 - don't use)
CACHE_EXPIRE=60  # seconds - ttl
MAX_CACHE_SIZE=1000000  # kb
MAX_CACHE_ITEM_SIZE=10000  #kb
TARGETS=http://example:8001;http://example2:8001
```
## Installation
### With docker  
```
$> docker pull freundallein/schooner
```
### With source
```
$> git clone git@github.com:freundallein/schooner.git
$> cd schooner
$> make build
```

## Usage
Docker-compose
```
version: "3.5"

networks:
  network:
    name: example-network
    driver: bridge

services:
  schooner:
    image: freundallein/schooner:latest
    container_name: schooner
    restart: always
    environment: 
      - PORT=8000
      - MACHINE_ID=0
      - USE_CACHE=1
      - CACHE_EXPIRE=60  # seconds
      - MAX_CACHE_SIZE=1000000  # kb
      - MAX_CACHE_ITEM_SIZE=10000  #kb
      - TARGETS=http://example:8001;http://example2:8001
      - STALE_TIMEOUT=1 # min
    networks: 
      - network
    ports:
      - 8000:8000
    depends_on: 
      - example
      - example2

  example:
    image: freundallein/schooner-example:latest
    container_name: example
    networks: 
      - network
  example2:
    image: freundallein/schooner-example:latest
    container_name: example2
    networks: 
      - network
```
## Metrics
Default prometheus metrics are available on `/schooner/metrics`  

## Healthcheck
Service healthcheck is avaliable on `/schooner/healthz`.  
