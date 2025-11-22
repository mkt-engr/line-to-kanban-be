package usecase

import "line-to-kanban-be/internal/domain/task"

type CreateTaskRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type TaskResponse struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

func ToTaskResponse(t *task.Task) *TaskResponse {
	return &TaskResponse{
		ID:      t.ID,
		UserID:  t.UserID,
		Content: t.Content,
		Status:  string(t.Status),
	}
}
