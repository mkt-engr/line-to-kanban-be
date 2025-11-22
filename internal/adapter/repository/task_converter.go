package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"line-to-kanban-be/internal/adapter/repository/db"
	"line-to-kanban-be/internal/domain/task"
)

// sqlcのdb.TaskからdomainのTaskに変換
func toTask(dbTask db.Task) *task.Task {
	// pgtype.UUIDをstring(UUID)に変換
	var idStr string
	if dbTask.ID.Valid {
		u, _ := uuid.FromBytes(dbTask.ID.Bytes[:])
		idStr = u.String()
	}

	return &task.Task{
		ID:        idStr,
		UserID:    dbTask.UserID,
		Content:   dbTask.Content,
		Status:    task.Status(dbTask.Status),
		CreatedAt: dbTask.CreatedAt.Time,
		UpdatedAt: dbTask.UpdatedAt.Time,
	}
}

// domainのTaskからsqlcのdb.TaskStatusに変換
func toDBStatus(status task.Status) db.TaskStatus {
	return db.TaskStatus(status)
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
