package repository

import (
	"context"

	"line-to-kanban-be/internal/adapter/repository/db"
	"line-to-kanban-be/internal/domain/task"
)

func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	_, err := r.queries.CreateTask(ctx, db.CreateTaskParams{
		Content: t.Content,
		Status:  toDBStatus(t.Status),
		UserID:  t.UserID,
	})
	return err
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id string, status task.Status) error {
	// ステータス更新機能は現在のsqlcに未実装のため、
	// 必要に応じてクエリを追加してください
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string, userID string) error {
	return r.queries.DeleteTask(ctx, db.DeleteTaskParams{
		ID:     toUUID(id),
		UserID: userID,
	})
}
