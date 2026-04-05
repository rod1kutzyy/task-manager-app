package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
)

func (r *repository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM notesapp.users
	WHERE id = $1;
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("query: %w", err)
	}

	userModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[UserModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d': %w",
				id, core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("collect user row: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
