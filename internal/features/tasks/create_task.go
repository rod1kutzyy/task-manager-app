package tasks_service

import (
	"context"
	"fmt"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
	tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"
)

func (s *service) CreateTask(ctx context.Context, in tasks_ports_in.CreateTaskParams) (tasks_ports_in.CreateTaskResult, error) {
	task := domain.CreateTask(
		in.Title,
		in.Description,
		in.AuthorUserID,
	)

	if err := task.Validate(); err != nil {
		return tasks_ports_in.CreateTaskResult{}, fmt.Errorf("validate task domain: %w", err)
	}

	repoParams := tasks_ports_out_repository.NewSaveTaskParams(task)
	repoResult, err := s.tasksRepository.SaveTask(ctx, repoParams)
	if err != nil {
		return tasks_ports_in.CreateTaskResult{}, fmt.Errorf("save task in repository: %w", err)
	}

	return tasks_ports_in.NewCreateTaskResult(repoResult.Task), nil
}
