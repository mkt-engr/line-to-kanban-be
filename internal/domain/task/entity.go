package task

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID        string
	UserID    string
	Content   string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(userID, content string) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		Status:    StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
