package repository

import (
	"context"

	"line-to-kanban-be/internal/domain/message"
)

func (r *MessageRepository) FindByID(ctx context.Context, id string) (*message.Message, error) {
	// IDで取得する機能は現在のsqlcに未実装のため、
	// 必要に応じてクエリを追加してください
	return nil, nil
}

func (r *MessageRepository) FindByUserID(ctx context.Context, userID string) ([]*message.Message, error) {
	dbMessages, err := r.queries.ListMessagesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	messages := make([]*message.Message, len(dbMessages))
	for i, dbMsg := range dbMessages {
		messages[i] = toMessage(dbMsg)
	}

	return messages, nil
}
