package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getStatisticsResponse struct {
	TasksCreated      int      `json:"tasks_created"`
	TasksCompleted    int      `json:"tasks_completed"`
	CompletedRate     *float64 `json:"completed_rate"`
	AvgCompletionTime *string  `json:"avg_completion_time"`
}

func toDTOFromDomain(statistics domain.Statistics) getStatisticsResponse {
	var avgTime *string
	if statistics.AvgCompletionTime != nil {
		duration := statistics.AvgCompletionTime.String()
		avgTime = &duration
	}

	return getStatisticsResponse{
		TasksCreated:      statistics.TasksCreated,
		TasksCompleted:    statistics.TasksCompleted,
		CompletedRate:     statistics.CompletedRate,
		AvgCompletionTime: avgTime,
	}
}

func (h *handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	resp := toDTOFromDomain(statistics)

	respHandler.JSONResponse(resp, http.StatusOK)
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
