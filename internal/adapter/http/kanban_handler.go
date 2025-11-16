package http

import (
	"encoding/json"
	"net/http"
)

type KanbanHandler struct{}

func NewKanbanHandler() *KanbanHandler {
	return &KanbanHandler{}
}

func (h *KanbanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "kanban status updated",
	})
}
