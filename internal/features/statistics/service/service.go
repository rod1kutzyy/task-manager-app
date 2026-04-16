package statistics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type service struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewService(statisticsRepository StatisticsRepository) *service {
	return &service{
		statisticsRepository: statisticsRepository,
	}
}
