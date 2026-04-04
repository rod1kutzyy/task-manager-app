package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (r *repository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO notesapp.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`

	rows, err := r.pool.Query(ctx, query, user.FullName, user.PhoneNumber)
	if err != nil {
		return domain.User{}, fmt.Errorf("query error: %w", err)
	}

	userModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[UserModel])
	if err != nil {
		return domain.User{}, fmt.Errorf("collect user row error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
