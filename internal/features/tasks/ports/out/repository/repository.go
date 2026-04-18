package tasks_ports_out_repository

import "context"

type TasksRepository interface {
	SaveTask(ctx context.Context, in SaveTaskParams) (SaveTaskResult, error)
}
