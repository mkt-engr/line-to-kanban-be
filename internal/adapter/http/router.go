package http

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Hello World endpoint
	mux.Handle("/", NewHelloHandler())

	// Health check endpoint
	mux.Handle("/healthz", NewHealthHandler())

	// Kanban status update endpoint
	mux.Handle("/kanban/status", NewKanbanHandler())

	return mux
}
