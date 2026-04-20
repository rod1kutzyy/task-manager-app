package tasks_adapters_out_repository_postgres

import (
	"context"
	"fmt"

	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (r *repository) ListTasks(
	ctx context.Context,
	in tasks_ports_out_repository.ListTasksParams,
) (tasks_ports_out_repository.ListTasksResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM notesapp.tasks
	%s
	ORDER BY created_at DESC, id ASC
	LIMIT $1
	OFFSET $2;
	`

	args := []any{in.Limit, in.Offset}
	if in.UserID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id = $3")
		args = append(args, in.UserID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return tasks_ports_out_repository.ListTasksResult{}, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
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
			return tasks_ports_out_repository.ListTasksResult{}, fmt.Errorf("scan task: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return tasks_ports_out_repository.ListTasksResult{}, fmt.Errorf("new rows: %w", err)
	}

	taskDomains := modelsToDomains(taskModels)

	return tasks_ports_out_repository.NewListTasksResult(taskDomains), nil
}
