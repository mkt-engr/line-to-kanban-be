package task

import "context"

type Repository interface {
	Save(ctx context.Context, t *Task) error
	FindByID(ctx context.Context, id string) (*Task, error)
	FindByUserID(ctx context.Context, userID string) ([]*Task, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
	Delete(ctx context.Context, id string, userID string) error
}
