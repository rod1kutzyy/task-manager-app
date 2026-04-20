package tasks_adapters_out_repository_cached

import (
	"context"
	"fmt"

	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *cachedRepository) DeleteTask(
	ctx context.Context,
	in tasks_ports_out_repository.DeleteTaskParams,
) (tasks_ports_out_repository.DeleteTaskResult, error) {
	task, ok := r.getTaskFromCache(ctx, in.ID)
	if !ok {
		mainRepoGetTaskResult, err := r.mainRepository.GetTask(
			ctx,
			tasks_ports_out_repository.NewGetTaskParams(in.ID),
		)
		if err != nil {
			return tasks_ports_out_repository.DeleteTaskResult{}, fmt.Errorf("get task: %w", err)
		}

		task = mainRepoGetTaskResult.Task
	}

	if _, err := r.mainRepository.DeleteTask(ctx, in); err != nil {
		return tasks_ports_out_repository.DeleteTaskResult{}, err
	}

	r.invalidateTask(ctx, task.AuthorUserID, &in.ID)

	return tasks_ports_out_repository.NewDeleteTaskResult(), nil
}
