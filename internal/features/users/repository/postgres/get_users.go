package users_postgres_repository

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
