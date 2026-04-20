package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) ListTasks(ctx context.Context, in tasks_ports_in.ListTasksParams) (tasks_ports_in.ListTasksResult, error) {
	if in.Limit != nil && *in.Limit < 0 {
		return tasks_ports_in.ListTasksResult{}, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if in.Offset != nil && *in.Offset < 0 {
		return tasks_ports_in.ListTasksResult{}, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	repoParams := tasks_ports_out_repository.NewListTasksParams(
		in.UserID,
		in.Limit,
		in.Offset,
	)
	repoResult, err := s.tasksRepository.ListTasks(ctx, repoParams)
	if err != nil {
		return tasks_ports_in.ListTasksResult{}, fmt.Errorf("list tasks from repository: %w", err)
	}

	return tasks_ports_in.NewListTasksResult(repoResult.Tasks), nil
}
