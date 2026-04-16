package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (s *service) GetTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
