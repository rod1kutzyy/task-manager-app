package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
)

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM notesapp.users
	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
