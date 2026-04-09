package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
)

type handler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

func NewHandler(tasksService TasksService) *handler {
	return &handler{
		tasksService: tasksService,
	}
}

func (h *handler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
	}
}
