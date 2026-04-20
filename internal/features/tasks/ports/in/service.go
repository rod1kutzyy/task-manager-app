package tasks_ports_in

import "context"

type TasksService interface {
	CreateTask(ctx context.Context, in CreateTaskParams) (CreateTaskResult, error)
	GetTask(ctx context.Context, in GetTaskParams) (GetTaskResult, error)
	GetTasks(ctx context.Context, in GetTasksParams) (GetTasksResult, error)
	DeleteTask(ctx context.Context, in DeleteTaskParams) (DeleteTaskResult, error)
}
