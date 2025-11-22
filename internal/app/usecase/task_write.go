package usecase

import (
	"context"

	"line-to-kanban-be/internal/domain/task"
)

func (u *TaskUsecase) CreateTask(ctx context.Context, req *CreateTaskRequest) (*TaskResponse, error) {
	t := task.NewTask(req.UserID, req.Content)

	if err := u.repo.Save(ctx, t); err != nil {
		return nil, err
	}

	return ToTaskResponse(t), nil
}

func (u *TaskUsecase) UpdateTaskStatus(ctx context.Context, id string, req *UpdateStatusRequest) error {
	status := task.Status(req.Status)
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id string, userID string) error {
	return u.repo.Delete(ctx, id, userID)
}
