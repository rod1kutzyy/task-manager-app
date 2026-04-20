package tasks_adapters_out_repository_cached

import (
	"encoding/json"
	"fmt"
	"strconv"
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

type TaksListModel []TaskModel

func (m *TaksListModel) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task list: %w", err)
	}

	return bytes, nil
}

func (m *TaksListModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize task list: %w", err)
	}

	return nil
}

func domainsToModels(tasks []domain.Task) TaksListModel {
	models := make(TaksListModel, len(tasks))

	for i, task := range tasks {
		models[i] = domainToModel(task)
	}

	return models
}

func modelsToDomains(list TaksListModel) []domain.Task {
	tasks := make([]domain.Task, len(list))

	for i, model := range list {
		tasks[i] = modelToDomain(model)
	}

	return tasks
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

func tasksListField(limit *int, offset *int) string {
	ptrStr := func(v *int) string {
		if v == nil {
			return "nil"
		}

		return strconv.Itoa(*v)
	}

	return fmt.Sprintf("%s:%s", ptrStr(limit), ptrStr(offset))
}
