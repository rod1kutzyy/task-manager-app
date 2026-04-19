package tasks_adapters_in_transport_http

import (
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
)

type handler struct {
	tasksService tasks_ports_in.TasksService
}

func NewHandler(tasksService tasks_ports_in.TasksService) *handler {
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
