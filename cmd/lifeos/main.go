package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("haiii"))
	})

	log.Println("Starting server at :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}

