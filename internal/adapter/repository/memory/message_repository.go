package memory

import (
	"context"
	"errors"
	"fmt"
	"line-to-kanban-be/internal/domain/message"
	"sync"
	"time"
)

type MessageRepository struct {
	mu       sync.RWMutex
	messages map[string]*message.Message
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		messages: make(map[string]*message.Message),
	}
}

func (r *MessageRepository) Save(ctx context.Context, msg *message.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if msg.ID == "" {
		msg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	r.messages[msg.ID] = msg
	return nil
}

func (r *MessageRepository) FindByID(ctx context.Context, id string) (*message.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msg, ok := r.messages[id]
	if !ok {
		return nil, errors.New("message not found")
	}

	return msg, nil
}

func (r *MessageRepository) FindAll(ctx context.Context) ([]*message.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messages := make([]*message.Message, 0, len(r.messages))
	for _, msg := range r.messages {
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepository) UpdateStatus(ctx context.Context, id string, status message.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[id]
	if !ok {
		return errors.New("message not found")
	}

	msg.UpdateStatus(status)
	return nil
}
