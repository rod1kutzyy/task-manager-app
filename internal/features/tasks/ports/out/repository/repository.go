package tasks_ports_out_repository

import "context"

type TasksRepository interface {
	SaveTask(ctx context.Context, in SaveTaskParams) (SaveTaskResult, error)
	GetTask(ctx context.Context, in GetTaskParams) (GetTaskResult, error)
	ListTasks(ctx context.Context, in ListTasksParams) (ListTasksResult, error)
	DeleteTask(ctx context.Context, in DeleteTaskParams) (DeleteTaskResult, error)
	UpdateTask(ctx context.Context, in UpdateTaskParams) (UpdateTaskResult, error)
}
