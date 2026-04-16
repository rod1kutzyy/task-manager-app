package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
)

func (r *repository) GetTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM notesapp.tasks
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%s': %w", id, core_errors.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("scan task: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
