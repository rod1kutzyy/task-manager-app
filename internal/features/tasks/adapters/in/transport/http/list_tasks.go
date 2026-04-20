package tasks_adapters_in_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
)

type getTasksResponse []taskDTOResponse

// ListTasks godoc
// @Summary List tasks
// @Description Returns tasks with optional filtering by `user_id` (UUID) and optional pagination.
// @Tags tasks
// @Produce json
// @Param user_id query string false "Filter tasks by author user ID (UUID)" Format(uuid)
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {array} taskDTOResponse "Tasks list"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get userID/limit/offset query params")
		return
	}

	serviceParams := tasks_ports_in.NewListTasksParams(userID, limit, offset)
	serviceResult, err := h.tasksService.ListTasks(ctx, serviceParams)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to list tasks")
		return
	}

	resp := getTasksResponse(taskDTOsFromDomains(serviceResult.Tasks))

	respHandler.JSONResponse(resp, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*uuid.UUID, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := request.GetUUIDQueryParam(r, userIDQueryParamKey)
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
