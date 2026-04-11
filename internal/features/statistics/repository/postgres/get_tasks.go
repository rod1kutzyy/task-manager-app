package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (r *repository) GetTasks(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM notesapp.tasks
	`)

	args := []any{}
	conditions := []string{}

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id;")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
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
