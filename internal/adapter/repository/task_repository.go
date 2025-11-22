package repository

import (
	"line-to-kanban-be/internal/adapter/repository/db"
)

type TaskRepository struct {
	queries db.Querier
}

func NewTaskRepository(queries db.Querier) *TaskRepository {
	return &TaskRepository{
		queries: queries,
	}
}
