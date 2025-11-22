package repository

import (
	"context"

	"line-to-kanban-be/internal/adapter/repository/db"
	"line-to-kanban-be/internal/domain/message"
)

func (r *MessageRepository) Save(ctx context.Context, msg *message.Message) error {
	_, err := r.queries.CreateMessage(ctx, db.CreateMessageParams{
		Content: msg.Content,
		Status:  toDBStatus(msg.Status),
		UserID:  msg.UserID,
	})
	return err
}

func (r *MessageRepository) UpdateStatus(ctx context.Context, id string, status message.Status) error {
	// ステータス更新機能は現在のsqlcに未実装のため、
	// 必要に応じてクエリを追加してください
	return nil
}

func (r *MessageRepository) Delete(ctx context.Context, id string, userID string) error {
	return r.queries.DeleteMessage(ctx, db.DeleteMessageParams{
		ID:     toUUID(id),
		UserID: userID,
	})
}
