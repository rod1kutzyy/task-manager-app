package tasks_transport_http

import (
	"time"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type taskDTOResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
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
