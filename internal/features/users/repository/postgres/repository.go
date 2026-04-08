package users_postgres_repository

import core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"

type repository struct {
	pool core_postgres_pool.Pool
}

func NewRepository(pool core_postgres_pool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
