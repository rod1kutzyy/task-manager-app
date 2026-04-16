package tasks_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type service struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *uuid.UUID, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (domain.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
	PatchTask(ctx context.Context, id uuid.UUID, task domain.Task) (domain.Task, error)
}

func NewService(tasksRepository TasksRepository) *service {
	return &service{
		tasksRepository: tasksRepository,
	}
}
