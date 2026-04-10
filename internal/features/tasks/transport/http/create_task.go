package tasks_transport_http

import (
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type createTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

type createTaskResponse taskDTOResponse

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	var req createTaskRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		respHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskDomain := domain.NewTaskUninitialized(req.Title, req.Description, req.AuthorUserID)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to create task")
		return
	}

	resp := createTaskResponse(taskDTOFromDomain(taskDomain))

	respHandler.JSONResponse(resp, http.StatusCreated)
}
