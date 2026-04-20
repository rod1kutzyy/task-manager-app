package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
)

type Pool struct {
	pool             *pgxpool.Pool
	operationTimeout time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		pool:             pool,
		operationTimeout: config.Timeout,
	}, nil
}

func (p *Pool) OperationTimeout() time.Duration {
	return p.operationTimeout
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	row := p.pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := p.pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}
