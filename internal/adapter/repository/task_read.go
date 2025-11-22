package repository

import (
	"context"

	"line-to-kanban-be/internal/domain/task"
)

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*task.Task, error) {
	// IDで取得する機能は現在のsqlcに未実装のため、
	// 必要に応じてクエリを追加してください
	return nil, nil
}

func (r *TaskRepository) FindByUserID(ctx context.Context, userID string) ([]*task.Task, error) {
	dbTasks, err := r.queries.ListTasksByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = toTask(dbTask)
	}

	return tasks, nil
}
