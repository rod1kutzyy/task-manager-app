package tasks_ports_out_repository

import (
	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type ListTasksParams struct {
	UserID *uuid.UUID
	Limit  *int
	Offset *int
}

func NewListTasksParams(userID *uuid.UUID, limit *int, offset *int) ListTasksParams {
	return ListTasksParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
}

type ListTasksResult struct {
	Tasks []domain.Task
}

func NewListTasksResult(tasks []domain.Task) ListTasksResult {
	return ListTasksResult{
		Tasks: tasks,
	}
}
