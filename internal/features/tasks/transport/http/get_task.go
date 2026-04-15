package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getTaskResponse taskDTOResponse

// GetTask godoc
// @Summary Get a task by ID
// @Description Returns a single task by numeric identifier.
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} getTaskResponse "Task found"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *handler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	taskID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get task")
		return
	}

	resp := getTaskResponse(taskDTOFromDomain(taskDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}
