package main

import (
	"io"
	"log"
	"net/http"

	"rate-limiter/limiter"
)

func main() {
	baseHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello from base handler")
	}

	http.HandleFunc("/", baseHandler)
	http.Handle("/limit", limiter.New())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
