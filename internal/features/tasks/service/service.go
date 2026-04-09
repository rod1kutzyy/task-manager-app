package tasks_service

import (
	"context"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type service struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
}

func NewService(tasksRepository TasksRepository) *service {
	return &service{
		tasksRepository: tasksRepository,
	}
}
