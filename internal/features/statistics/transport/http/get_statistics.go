package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getStatisticsResponse struct {
	TasksCreated      int      `json:"tasks_created" example:"50"`
	TasksCompleted    int      `json:"tasks_completed" example:"10"`
	CompletedRate     *float64 `json:"completed_rate" example:"20"`
	AvgCompletionTime *string  `json:"avg_completion_time" example:"1m30s"`
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

// GetStatistics godoc
// @Summary Get task statistics
// @Description Returns task statistics with optional filtering by `user_id` (UUID) and/or date range.
// @Tags statistics
// @Produce json
// @Param user_id query string false "Filter statistics by user ID (UUID)" Format(uuid)
// @Param from query string false "Start date (inclusive), format: YYYY-MM-DD"
// @Param to query string false "End date (exclusive), format: YYYY-MM-DD"
// @Success 200 {object} getStatisticsResponse "Statistics response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /statistics [get]
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

func getUserIDFromToQueryParams(r *http.Request) (*uuid.UUID, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := request.GetUUIDQueryParam(r, userIDQueryParamKey)
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
