package usecase

import (
	"line-to-kanban-be/internal/domain/task"
)

type TaskUsecase struct {
	repo task.Repository
}

func NewTaskUsecase(repo task.Repository) *TaskUsecase {
	return &TaskUsecase{
		repo: repo,
	}
}
