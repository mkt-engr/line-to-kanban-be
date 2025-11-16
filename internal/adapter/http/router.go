package http

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Hello World endpoint
	mux.Handle("/", NewHelloHandler())

	return mux
}
