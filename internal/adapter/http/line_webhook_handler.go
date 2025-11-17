package http

import (
	"line-to-kanban-be/internal/adapter/line"
	"net/http"
)

type LineWebhookHandler struct {
	webhookHandler *line.WebhookHandler
}

func NewLineWebhookHandler(webhookHandler *line.WebhookHandler) *LineWebhookHandler {
	return &LineWebhookHandler{
		webhookHandler: webhookHandler,
	}
}

func (h *LineWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.webhookHandler.Handle(w, r)
}
