package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Message: "hello world",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	port := ":8080"
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
