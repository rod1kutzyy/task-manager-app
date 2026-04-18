package tasks_adapters_out_repository_postgres

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

func modelToDomain(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}
