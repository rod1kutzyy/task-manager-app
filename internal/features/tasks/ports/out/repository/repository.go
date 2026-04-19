package tasks_ports_out_repository

import "context"

type TasksRepository interface {
	SaveTask(ctx context.Context, in SaveTaskParams) (SaveTaskResult, error)
	GetTask(ctx context.Context, in GetTaskParams) (GetTaskResult, error)
	DeleteTask(ctx context.Context, in DeleteTaskParams) (DeleteTaskResult, error)
}
