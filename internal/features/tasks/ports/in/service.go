package tasks_ports_in

import "context"

type TasksService interface {
	CreateTask(ctx context.Context, in CreateTaskParams) (CreateTaskResult, error)
	GetTask(ctx context.Context, in GetTaskParams) (GetTaskResult, error)
	DeleteTask(ctx context.Context, in DeleteTaskParams) (DeleteTaskResult, error)
}
