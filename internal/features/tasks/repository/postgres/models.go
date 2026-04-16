package tasks_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type TaskModel struct {
	ID           uuid.UUID  `db:"id"`
	Version      int        `db:"version"`
	Title        string     `db:"title"`
	Description  *string    `db:"description"`
	Completed    bool       `db:"completed"`
	CreatedAt    time.Time  `db:"created_at"`
	CompletedAt  *time.Time `db:"completed_at"`
	AuthorUserID uuid.UUID  `db:"author_user_id"`
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
