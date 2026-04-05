package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (r *repository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM notesapp.users
	ORDER BY id
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	userModels, err := pgx.CollectRows(rows, pgx.RowToStructByPos[UserModel])
	if err != nil {
		return nil, fmt.Errorf("collect users rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
