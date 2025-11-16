package http

import (
	"encoding/json"
	"line-to-kanban-be/internal/app/message"
	"net/http"
)

type KanbanHandler struct {
	usecase *message.Usecase
}

func NewKanbanHandler(usecase *message.Usecase) *KanbanHandler {
	return &KanbanHandler{
		usecase: usecase,
	}
}

func (h *KanbanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "kanban status updated",
	})
}
