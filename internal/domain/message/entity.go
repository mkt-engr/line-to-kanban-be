package message

import "time"

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Message struct {
	ID        string
	Content   string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMessage(content string) *Message {
	now := time.Now()
	return &Message{
		Content:   content,
		Status:    StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *Message) UpdateStatus(status Status) {
	m.Status = status
	m.UpdatedAt = time.Now()
}
