package tasks_ports_out_repository

import "github.com/google/uuid"

type DeleteTaskParams struct {
	ID uuid.UUID
}

func NewDeleteTaskParams(id uuid.UUID) DeleteTaskParams {
	return DeleteTaskParams{
		ID: id,
	}
}

type DeleteTaskResult struct{}

func NewDeleteTaskResult() DeleteTaskResult {
	return DeleteTaskResult{}
}
