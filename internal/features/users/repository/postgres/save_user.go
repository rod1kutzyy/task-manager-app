package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (r *repository) SaveUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO notesapp.users (id, version, full_name, phone_number)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Version,
		user.FullName,
		user.PhoneNumber,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
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
