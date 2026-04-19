package tasks_adapters_out_repository_cached

import (
	"context"

	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *cachedRepository) GetTask(
	ctx context.Context,
	in tasks_ports_out_repository.GetTaskParams,
) (tasks_ports_out_repository.GetTaskResult, error) {
	if task, ok := r.getTaskFromCache(ctx, in.ID); ok {
		return tasks_ports_out_repository.NewGetTaskResult(task), nil
	}

	mainRepoResult, err := r.mainRepository.GetTask(ctx, in)
	if err != nil {
		return tasks_ports_out_repository.GetTaskResult{}, err
	}

	r.cacheTask(ctx, mainRepoResult.Task)

	return mainRepoResult, nil
}
