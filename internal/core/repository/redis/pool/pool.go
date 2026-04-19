package core_redis_pool

import (
	"context"
	"time"
)

type Pool interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, ttl time.Duration) StatusCmd
	Del(ctx context.Context, keys ...string) IntCmd
	Close() error

	TTL() time.Duration
}

type StringCmd interface {
	Bytes() ([]byte, error)
}

type StatusCmd interface {
	Err() error
}

type IntCmd interface {
	Err() error
}
