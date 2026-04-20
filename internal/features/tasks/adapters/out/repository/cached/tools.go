package tasks_adapters_out_repository_cached

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	core_redis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool"
	"go.uber.org/zap"
)

func (c *cachedRepository) getTaskFromCache(ctx context.Context, id uuid.UUID) (domain.Task, bool) {
	logger := core_logger.FromContext(ctx)

	key := taskKey(id)

	bytes, err := c.pool.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.ErrNotFound) {
			logger.Error("read from cache", zap.Error(err))
		}

		return domain.Task{}, false
	}

	var taskModel TaskModel
	if err := taskModel.Deserialize(bytes); err != nil {
		logger.Error("deserialize cached task", zap.Error(err))
		return domain.Task{}, false
	}

	taskDomain := modelToDomain(taskModel)

	return taskDomain, true
}

func (r *cachedRepository) cacheTask(ctx context.Context, task domain.Task) bool {
	logger := core_logger.FromContext(ctx)

	taskModel := domainToModel(task)
	bytes, err := taskModel.Serialize()
	if err != nil {
		logger.Error("serialize task", zap.Error(err))
		return false
	} else {
		if err := r.pool.Set(ctx, taskKey(taskModel.ID), bytes, r.pool.TTL()).Err(); err != nil {
			logger.Error("set task in cache", zap.Error(err))
			return false
		}
	}

	return true
}

func (r *cachedRepository) invalidateTask(ctx context.Context, userID uuid.UUID, taskID *uuid.UUID) {
	logger := core_logger.FromContext(ctx)

	invalidateKeys := []string{
		tasksListKey(nil),
		tasksListKey(&userID),
	}

	if taskID != nil {
		invalidateKeys = append(invalidateKeys, taskKey(*taskID))
	}

	if err := r.pool.Del(ctx, invalidateKeys...).Err(); err != nil {
		logger.Error("invalidate cached tasks lists", zap.Error(err))
	}
}
