package tasks_adapters_out_repository_cached

import (
	"context"

	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *cachedRepository) UpdateTask(
	ctx context.Context,
	in tasks_ports_out_repository.UpdateTaskParams,
) (tasks_ports_out_repository.UpdateTaskResult, error) {
	mainRepoResult, err := r.mainRepository.UpdateTask(ctx, in)
	if err != nil {
		return tasks_ports_out_repository.UpdateTaskResult{}, err
	}

	task := mainRepoResult.Task

	if ok := r.cacheTask(ctx, task); !ok {
		r.invalidateTask(ctx, task.AuthorUserID, &task.ID)
	} else {
		r.invalidateTask(ctx, task.AuthorUserID, nil)
	}

	return mainRepoResult, nil
}
