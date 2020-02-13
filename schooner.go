package main

import (
	"fmt"
	"github.com/freundallein/schooner/httpserv"
	"log"
	"os"
	"time"
)

const (
	timeFormat = "02.01.2006 15:04:05"

	portKey     = "PORT"
	defaultPort = "8000"
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

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	port, err := getEnv(portKey, defaultPort)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	options := &httpserv.Options{
		Port: port,
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
