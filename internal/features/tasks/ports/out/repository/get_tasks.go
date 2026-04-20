package tasks_ports_out_repository

import (
	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type GetTasksParams struct {
	UserID *uuid.UUID
	Limit  *int
	Offset *int
}

func NewGetTasksParams(userID *uuid.UUID, limit *int, offset *int) GetTasksParams {
	return GetTasksParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
}

type GetTasksResult struct {
	Tasks []domain.Task
}

func NewGetTasksResult(tasks []domain.Task) GetTasksResult {
	return GetTasksResult{
		Tasks: tasks,
	}
}
