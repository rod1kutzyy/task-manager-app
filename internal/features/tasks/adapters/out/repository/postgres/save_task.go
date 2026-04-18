package tasks_adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	core_postgres_pool "github.com/rod1kutzyy/task-manager-app/internal/core/repository/postgres/pool"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *repository) SaveTask(
	ctx context.Context,
	params tasks_ports_out_repository.SaveTaskParams,
) (tasks_ports_out_repository.SaveTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO notesapp.tasks (id, version, title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	task := params.Task
	row := r.pool.QueryRow(
		ctx,
		query,
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
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
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return tasks_ports_out_repository.SaveTaskResult{}, fmt.Errorf(
				"%v: user with id='%s': %w",
				err, task.AuthorUserID, core_errors.ErrNotFound,
			)
		}

		return tasks_ports_out_repository.SaveTaskResult{}, fmt.Errorf("scan task: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return tasks_ports_out_repository.NewSaveTaskResult(taskDomain), nil
}
