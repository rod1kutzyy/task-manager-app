package tasks_ports_in

import "context"

type TasksService interface {
	CreateTask(ctx context.Context, in CreateTaskParams) (CreateTaskResult, error)
}
