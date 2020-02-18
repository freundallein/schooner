package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnv(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return fallback, nil
}
func main() {
	port, err := getEnv("PORT", "8001")
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	http.HandleFunc("/fast", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HELLO"))
	})
	http.HandleFunc("/nocache", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Pragma", "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("NO CACHE HELLO"))
	})
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SLOW HELLO"))
	})
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./files"))))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
