package tasks_adapters_out_repository_cached

import (
	"context"
	"errors"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	core_redis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
	"go.uber.org/zap"
)

func (r *cachedRepository) ListTasks(
	ctx context.Context,
	in tasks_ports_out_repository.ListTasksParams,
) (tasks_ports_out_repository.ListTasksResult, error) {
	logger := core_logger.FromContext(ctx)

	key := tasksListKey(in.UserID)
	filed := tasksListField(in.Limit, in.Offset)

	bytes, err := r.pool.HGet(ctx, key, filed).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.ErrNotFound) {
			logger.Error("hget task list from cache", zap.Error(err))
		}
	} else {
		var taskListModel TaksListModel
		if err := taskListModel.Deserialize(bytes); err != nil {
			logger.Error("deserialize cached task list", zap.Error(err))
		} else {
			tasks := modelsToDomains(taskListModel)

			return tasks_ports_out_repository.NewListTasksResult(tasks), nil
		}
	}

	mainRepoResult, err := r.mainRepository.ListTasks(ctx, in)
	if err != nil {
		return tasks_ports_out_repository.ListTasksResult{}, err
	}

	taskListModel := domainsToModels(mainRepoResult.Tasks)

	bytes, err = taskListModel.Serialize()
	if err != nil {
		logger.Error("serialize task list", zap.Error(err))
	} else {
		if err := r.pool.HSet(ctx, key, filed, bytes).Err(); err != nil {
			logger.Error("hset task list in cache")
		}
	}

	return mainRepoResult, nil
}
