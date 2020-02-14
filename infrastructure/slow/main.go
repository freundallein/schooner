package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		time.Sleep(8 * time.Second)
		w.Header().Add("KEY", "VALUE")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SLOW HELLO"))
	})
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		panic(err)
	}
}
