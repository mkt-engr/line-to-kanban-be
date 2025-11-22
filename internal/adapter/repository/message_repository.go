package repository

import (
	"line-to-kanban-be/internal/adapter/repository/db"
)

type MessageRepository struct {
	queries db.Querier
}

func NewMessageRepository(queries db.Querier) *MessageRepository {
	return &MessageRepository{
		queries: queries,
	}
}
