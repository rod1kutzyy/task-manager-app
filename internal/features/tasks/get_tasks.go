package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) GetTasks(ctx context.Context, in tasks_ports_in.GetTasksParams) (tasks_ports_in.GetTasksResult, error) {
	if in.Limit != nil && *in.Limit < 0 {
		return tasks_ports_in.GetTasksResult{}, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if in.Offset != nil && *in.Offset < 0 {
		return tasks_ports_in.GetTasksResult{}, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	repoParams := tasks_ports_out_repository.NewGetTasksParams(
		in.UserID,
		in.Limit,
		in.Offset,
	)
	repoResult, err := s.tasksRepository.GetTasks(ctx, repoParams)
	if err != nil {
		return tasks_ports_in.GetTasksResult{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks_ports_in.NewGetTasksResult(repoResult.Tasks), nil
}
