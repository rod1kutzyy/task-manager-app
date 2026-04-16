package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
)

func (r *repository) PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE notesapp.users
	SET
		full_name = $1,
		phone_number = $2,
		version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, id, user.Version)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%s' concurrently accessed: %w",
				id, core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
