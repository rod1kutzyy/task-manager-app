package tasks_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type taskDTOResponse struct {
	ID           uuid.UUID  `json:"id" example:"b5fc5eb1-65b6-4a28-ac34-612fd14cd39c"`
	Version      int        `json:"version" example:"2"`
	Title        string     `json:"title" example:"Walk the dog"`
	Description  *string    `json:"description" example:"Morning walk at 06:30"`
	Completed    bool       `json:"completed" example:"true"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-04-15T08:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"2026-04-15T09:00:00Z"`
	AuthorUserID uuid.UUID  `json:"author_user_id" example:"279ee73e-0132-4d5e-9498-cfe6fb43742c"`
}

func taskDTOFromDomain(task domain.Task) taskDTOResponse {
	return taskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []taskDTOResponse {
	dtos := make([]taskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
