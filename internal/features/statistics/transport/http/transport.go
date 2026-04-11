package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
)

type handler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(ctx context.Context, userID *int, from *time.Time, to *time.Time) (domain.Statistics, error)
}

func NewHandler(statisticsService StatisticsService) *handler {
	return &handler{
		statisticsService: statisticsService,
	}
}

func (h *handler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
