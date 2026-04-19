package tasks_service

import (
	"context"
	"fmt"

	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) DeleteTask(ctx context.Context, in tasks_ports_in.DeleteTaskParams) (tasks_ports_in.DeleteTaskResult, error) {
	repoParams := tasks_ports_out_repository.NewDeleteTaskParams(in.ID)
	if _, err := s.tasksRepository.DeleteTask(ctx, repoParams); err != nil {
		return tasks_ports_in.DeleteTaskResult{}, fmt.Errorf("delete task from repository: %w", err)
	}

	return tasks_ports_in.DeleteTaskResult{}, nil
}
