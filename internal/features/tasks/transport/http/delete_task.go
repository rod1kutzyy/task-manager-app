package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

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
