package core_goredis_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	core_redis_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/redis/pool"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPool(ctx context.Context, cfg Config) (*Pool, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) core_redis_pool.StatusCmd {
	cmd := p.client.Set(ctx, key, value, ttl)

	return goredisStatusCmd{cmd}
}

func (p *Pool) Del(ctx context.Context, keys ...string) core_redis_pool.IntCmd {
	cmd := p.client.Del(ctx, keys...)

	return goredisIntCmd{cmd}
}

func (p *Pool) Close() error {
	return p.client.Close()
}

func (p *Pool) TTL() time.Duration {
	return p.ttl
}
