package http

import (
	"line-to-kanban-be/internal/adapter/line"
	"net/http"
)

func NewRouter(lineWebhookHandler *line.WebhookHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// LINE Webhook endpoint
	mux.Handle("/webhook/line", NewLineWebhookHandler(lineWebhookHandler))

	return mux
}
