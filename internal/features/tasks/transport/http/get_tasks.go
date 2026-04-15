package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getTasksResponse []taskDTOResponse

// GetTasks godoc
// @Summary List tasks
// @Description Returns tasks with optional filtering by `user_id` and optional pagination.
// @Tags tasks
// @Produce json
// @Param user_id query int false "Filter tasks by author user ID"
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {array} taskDTOResponse "Tasks list"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	tasksDomain, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	resp := getTasksResponse(taskDTOsFromDomains(tasksDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
