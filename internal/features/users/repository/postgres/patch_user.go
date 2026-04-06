package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
)

func (r *repository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE notesapp.users
	SET
		full_name = $1,
		phone_number = $2,
		version = version+1
	WHERE id = $3 AND version = $4
	RETURNING id, version, full_name, phone_number;
	`

	rows, err := r.pool.Query(ctx, query, user.FullName, user.PhoneNumber, id, user.Version)
	if err != nil {
		return domain.User{}, fmt.Errorf("query: %w", err)
	}

	userModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[UserModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				id, core_errors.ErrConflict,
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
