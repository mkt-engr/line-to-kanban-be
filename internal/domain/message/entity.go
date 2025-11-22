package message

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

type Message struct {
	ID        string
	UserID    string
	Content   string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMessage(userID, content string) *Message {
	now := time.Now()
	return &Message{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		Status:    StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
