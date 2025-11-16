package http

import (
	"encoding/json"
	"line-to-kanban-be/internal/app/message"
	"net/http"
)

type LineWebhookHandler struct {
	usecase *message.Usecase
}

func NewLineWebhookHandler(usecase *message.Usecase) *LineWebhookHandler {
	return &LineWebhookHandler{
		usecase: usecase,
	}
}

func (h *LineWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "LINE webhook received",
	})
}
