package tasks_adapters_out_repository_cached

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type TaskModel struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID uuid.UUID  `json:"author_user_id"`
}

func (m *TaskModel) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task: %w", err)
	}

	return bytes, nil
}

func (m *TaskModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize tasK: %w", err)
	}

	return nil
}

func domainToModel(task domain.Task) TaskModel {
	return TaskModel{
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

func taskKey(id uuid.UUID) string {
	return fmt.Sprintf("task:%s", id)
}

func tasksListKey(userID *uuid.UUID) string {
	if userID == nil {
		return "tasks:all"
	}

	return fmt.Sprintf("tasks:%s", *userID)
}
