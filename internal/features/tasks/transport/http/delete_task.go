package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task by numeric identifier.
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204 "Task deleted"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	taskID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		respHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	respHandler.NoContentResponse()
}
