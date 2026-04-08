package tasks_postgres_repository

import (
	"time"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	tasksDomain := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		tasksDomain[i] = domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.CreatedAt,
			task.CompletedAt,
			task.AuthorUserID,
		)
	}

	return tasksDomain
}
