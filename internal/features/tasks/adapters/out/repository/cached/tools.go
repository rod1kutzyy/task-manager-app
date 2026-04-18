package tasks_adapters_out_repository_cached

import (
	"context"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"go.uber.org/zap"
)

func (r *cachedRepository) cacheTask(ctx context.Context, task domain.Task) {
	logger := core_logger.FromContext(ctx)

	taskModel := domainToModel(task)
	bytes, err := taskModel.Serialize()
	if err != nil {
		logger.Error("serialize task", zap.Error(err))
	} else {
		if err := r.pool.Set(ctx, taskKey(taskModel.ID), bytes, r.pool.TTL()).Err(); err != nil {
			logger.Error("set task in cache", zap.Error(err))
		}
	}
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
