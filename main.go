package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	go func() {
		http.ListenAndServe(":10000", &handler{})
	}()

	http.ListenAndServe(":"+port, &handler{})
}

type handler struct {
}

func (m *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi"))
}
