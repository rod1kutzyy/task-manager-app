package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (r *repository) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM notesapp.tasks
	%s
	ORDER BY id
	LIMIT $1
	OFFSET $2;
	`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
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
			return nil, fmt.Errorf("scan task: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	tasksDomains := taskDomainsFromModels(taskModels)

	return tasksDomains, nil

}
