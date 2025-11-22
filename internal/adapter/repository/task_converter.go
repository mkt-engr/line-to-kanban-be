package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"line-to-kanban-be/internal/adapter/repository/db"
	"line-to-kanban-be/internal/domain/task"
)

// sqlcのdb.MessageからdomainのTaskに変換
func toTask(dbMsg db.Message) *task.Task {
	// pgtype.UUIDをstring(UUID)に変換
	var idStr string
	if dbMsg.ID.Valid {
		u, _ := uuid.FromBytes(dbMsg.ID.Bytes[:])
		idStr = u.String()
	}

	return &task.Task{
		ID:        idStr,
		UserID:    dbMsg.UserID,
		Content:   dbMsg.Content,
		Status:    task.Status(dbMsg.Status),
		CreatedAt: dbMsg.CreatedAt.Time,
		UpdatedAt: dbMsg.UpdatedAt.Time,
	}
}

// domainのTaskからsqlcのdb.MessageStatusに変換
func toDBStatus(status task.Status) db.MessageStatus {
	return db.MessageStatus(status)
}

// stringからpgtype.UUIDに変換
func toUUID(id string) pgtype.UUID {
	var uuid pgtype.UUID
	// UUIDパース
	if err := uuid.Scan(id); err != nil {
		// エラー処理
	}
	return uuid
}
