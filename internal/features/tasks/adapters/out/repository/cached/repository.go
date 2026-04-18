package tasks_adapters_out_repository_cached

import (
	core_redis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

type cachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository tasks_ports_out_repository.TasksRepository
}

func NewCachedRepository(pool core_redis_pool.Pool, mainRepository tasks_ports_out_repository.TasksRepository) *cachedRepository {
	return &cachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
	}
}
