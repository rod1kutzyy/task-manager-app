package tasks_adapters_out_repository_cached

import (
	"context"

	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *cachedRepository) SaveTask(
	ctx context.Context,
	in tasks_ports_out_repository.SaveTaskParams,
) (tasks_ports_out_repository.SaveTaskResult, error) {
	mainRepoResult, err := r.mainRepository.SaveTask(ctx, in)
	if err != nil {
		return tasks_ports_out_repository.SaveTaskResult{}, err
	}

	task := mainRepoResult.Task

	_ = r.cacheTask(ctx, task)
	r.invalidateTask(ctx, task.AuthorUserID, nil)

	return mainRepoResult, nil
}
