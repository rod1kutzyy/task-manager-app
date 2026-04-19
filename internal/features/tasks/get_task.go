package tasks_service

import (
	"context"
	"fmt"

	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) GetTask(ctx context.Context, in tasks_ports_in.GetTaskParams) (tasks_ports_in.GetTaskResult, error) {
	repoParams := tasks_ports_out_repository.NewGetTaskParams(in.ID)
	repoResult, err := s.tasksRepository.GetTask(ctx, repoParams)
	if err != nil {
		return tasks_ports_in.GetTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	return tasks_ports_in.NewGetTaskResult(repoResult.Task), nil
}
