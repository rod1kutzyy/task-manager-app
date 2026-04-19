package tasks_adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *repository) GetTask(
	ctx context.Context,
	in tasks_ports_out_repository.GetTaskParams,
) (tasks_ports_out_repository.GetTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM notesapp.tasks
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, in.ID)

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
			return tasks_ports_out_repository.GetTaskResult{}, fmt.Errorf(
				"task with id='%s': %w",
				in.ID, core_errors.ErrNotFound,
			)
		}
		return tasks_ports_out_repository.GetTaskResult{}, fmt.Errorf("scan task: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return tasks_ports_out_repository.NewGetTaskResult(taskDomain), nil
}
