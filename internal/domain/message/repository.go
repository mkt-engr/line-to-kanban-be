package message

import "context"

type Repository interface {
	Save(ctx context.Context, msg *Message) error
	FindByID(ctx context.Context, id string) (*Message, error)
	FindAll(ctx context.Context) ([]*Message, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
}
