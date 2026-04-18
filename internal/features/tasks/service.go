package tasks_service

import tasks_ports_out_repository "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/out/repository"

type service struct {
	tasksRepository tasks_ports_out_repository.TasksRepository
}

func NewService(tasksRepository tasks_ports_out_repository.TasksRepository) *service {
	return &service{
		tasksRepository: tasksRepository,
	}
}
