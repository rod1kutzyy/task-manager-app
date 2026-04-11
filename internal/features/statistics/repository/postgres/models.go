package statistics_postgres_repository

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

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	tasksDomain := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		tasksDomain[i] = taskDomainFromModel(task)
	}

	return tasksDomain
}
