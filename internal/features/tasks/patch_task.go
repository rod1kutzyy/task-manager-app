package tasks_service

import (
	"context"
	"fmt"

	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) PatchTask(ctx context.Context, in tasks_ports_in.PatchTaskParams) (tasks_ports_in.PatchTaskResult, error) {
	repoGetTaskParams := tasks_ports_out_repository.NewGetTaskParams(in.ID)
	repoGetTaskResult, err := s.tasksRepository.GetTask(ctx, repoGetTaskParams)
	if err != nil {
		return tasks_ports_in.PatchTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	task := repoGetTaskResult.Task
	if err := task.ApplyPatch(in.Patch); err != nil {
		return tasks_ports_in.PatchTaskResult{}, fmt.Errorf("apply task patch: %w", err)
	}

	repoUpdateTaskParams := tasks_ports_out_repository.NewUpdateTaskParams(task)
	repoUpdateTaskResult, err := s.tasksRepository.UpdateTask(ctx, repoUpdateTaskParams)
	if err != nil {
		return tasks_ports_in.PatchTaskResult{}, fmt.Errorf("update task in repository: %w", err)
	}

	return tasks_ports_in.NewPatchTaskResult(repoUpdateTaskResult.Task), nil
}
