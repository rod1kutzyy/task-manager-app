package core_goredis_pool

import (
	"errors"

	"github.com/redis/go-redis/v9"
	core_redis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool"
)

type goredisStringCmd struct {
	*redis.StringCmd
}

func (c goredisStringCmd) Bytes() ([]byte, error) {
	data, err := c.StringCmd.Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return data, nil
}

type goredisStatusCmd struct {
	*redis.StatusCmd
}

type goredisIntCmd struct {
	*redis.IntCmd
}

func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.ErrNotFound
	}

	return err
}
