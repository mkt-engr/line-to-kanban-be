package http

import (
	"line-to-kanban-be/internal/adapter/line"
	"net/http"
)

func NewRouter(lineWebhookHandler *line.WebhookHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Hello World endpoint
	mux.Handle("/", NewHelloHandler())

	// Health check endpoint
	mux.Handle("/healthz", NewHealthHandler())

	// Webhook endpoints
	mux.Handle("/webhook/line", NewLineWebhookHandler(lineWebhookHandler))

	// Kanban status update endpoint
	mux.Handle("/kanban/status", NewKanbanHandler())

	return mux
}
