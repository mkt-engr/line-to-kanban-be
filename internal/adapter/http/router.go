package http

import (
	"line-to-kanban-be/internal/app/message"
	"net/http"
)

func NewRouter(usecase *message.Usecase) *http.ServeMux {
	mux := http.NewServeMux()

	// Hello World endpoint
	mux.Handle("/", NewHelloHandler())

	// Health check
	mux.Handle("/healthz", NewHealthHandler())

	// LINE webhook
	mux.Handle("/line/webhook", NewLineWebhookHandler(usecase))

	// Admin API
	mux.Handle("/admin/tasks", NewAdminHandler(usecase))

	// Kanban status update
	mux.Handle("/kanban/status", NewKanbanHandler(usecase))

	return mux
}
