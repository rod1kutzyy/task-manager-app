package tasks_adapters_out_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/rod1kutzyy/task-manager-app/internal/core/errors"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *repository) DeleteTask(
	ctx context.Context,
	in tasks_ports_out_repository.DeleteTaskParams,
) (tasks_ports_out_repository.DeleteTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM notesapp.tasks
	WHERE id = $1;
	`

	tag, err := r.pool.Exec(ctx, query, in.ID)
	if err != nil {
		return tasks_ports_out_repository.DeleteTaskResult{}, fmt.Errorf("exec query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return tasks_ports_out_repository.DeleteTaskResult{}, fmt.Errorf(
			"task with id='%s': %w",
			in.ID, core_errors.ErrNotFound,
		)
	}

	return tasks_ports_out_repository.NewDeleteTaskResult(), nil
}
