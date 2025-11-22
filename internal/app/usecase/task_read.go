package usecase

import (
	"context"
)

func (u *TaskUsecase) ListTasksByUser(ctx context.Context, userID string) ([]*TaskResponse, error) {
	tasks, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = ToTaskResponse(task)
	}

	return responses, nil
}
