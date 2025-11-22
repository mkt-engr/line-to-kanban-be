package message

import "context"

type Repository interface {
	Save(ctx context.Context, msg *Message) error
	FindByID(ctx context.Context, id string) (*Message, error)
	FindByUserID(ctx context.Context, userID string) ([]*Message, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
	Delete(ctx context.Context, id string, userID string) error
}
