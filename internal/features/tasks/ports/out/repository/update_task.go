package tasks_ports_out_repository

import "github.com/rod1kutzyy/task-manager-app/internal/core/domain"

type UpdateTaskParams struct {
	Task domain.Task
}

func NewUpdateTaskParams(task domain.Task) UpdateTaskParams {
	return UpdateTaskParams{
		Task: task,
	}
}

type UpdateTaskResult struct {
	Task domain.Task
}

func NewUpdateTaskResult(task domain.Task) UpdateTaskResult {
	return UpdateTaskResult{
		Task: task,
	}
}
