package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HELLO"))
	})
	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		panic(err)
	}
}
