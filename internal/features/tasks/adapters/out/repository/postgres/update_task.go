package tasks_adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *repository) UpdateTask(
	ctx context.Context,
	in tasks_ports_out_repository.UpdateTaskParams,
) (tasks_ports_out_repository.UpdateTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE notesapp.tasks
	SET
		title = $1,
		description = $2,
		completed = $3,
		completed_at = $4,
		version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		in.Task.Title,
		in.Task.Description,
		in.Task.Completed,
		in.Task.CompletedAt,
		in.Task.ID,
		in.Task.Version,
	)

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
			return tasks_ports_out_repository.UpdateTaskResult{}, fmt.Errorf(
				"task with id='%s' concurrently accessed: %w",
				in.Task.ID, core_errors.ErrConflict,
			)
		}

		return tasks_ports_out_repository.UpdateTaskResult{}, fmt.Errorf("scan task: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return tasks_ports_out_repository.NewUpdateTaskResult(taskDomain), nil
}
