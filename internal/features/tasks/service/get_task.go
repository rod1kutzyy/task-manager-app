package tasks_service

import (
	"context"
	"fmt"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (s *service) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
