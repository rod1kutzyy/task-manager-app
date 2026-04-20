package tasks_ports_in

import "context"

type TasksService interface {
	CreateTask(ctx context.Context, in CreateTaskParams) (CreateTaskResult, error)
	GetTask(ctx context.Context, in GetTaskParams) (GetTaskResult, error)
	ListTasks(ctx context.Context, in ListTasksParams) (ListTasksResult, error)
	DeleteTask(ctx context.Context, in DeleteTaskParams) (DeleteTaskResult, error)
	PatchTask(ctx context.Context, in PatchTaskParams) (PatchTaskResult, error)
}
