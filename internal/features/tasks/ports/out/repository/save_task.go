package tasks_ports_out_repository

import "github.com/rod1kutzyy/task-manager-app/internal/core/domain"

type SaveTaskParams struct {
	Task domain.Task
}

func NewSaveTaskParams(task domain.Task) SaveTaskParams {
	return SaveTaskParams{
		Task: task,
	}
}

type SaveTaskResult struct {
	Task domain.Task
}

func NewSaveTaskResult(task domain.Task) SaveTaskResult {
	return SaveTaskResult{
		Task: task,
	}
}
