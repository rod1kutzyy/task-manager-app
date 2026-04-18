package core_goredis_pool

import "github.com/redis/go-redis/v9"

type goredisStatusCmd struct {
	*redis.StatusCmd
}

type goredisIntCmd struct {
	*redis.IntCmd
}
