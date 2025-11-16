package http

import (
	"encoding/json"
	"line-to-kanban-be/internal/app/message"
	"net/http"
)

type AdminHandler struct {
	usecase *message.Usecase
}

func NewAdminHandler(usecase *message.Usecase) *AdminHandler {
	return &AdminHandler{
		usecase: usecase,
	}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "admin task created",
	})
}
