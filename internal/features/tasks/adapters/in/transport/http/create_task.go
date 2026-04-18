package tasks_adapters_in_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
	tasks_ports_in "github.com/rod1kutzyy/task-manager-app/internal/features/tasks/ports/in"
)

type createTaskRequest struct {
	Title        string    `json:"title" validate:"required,min=1,max=100" example:"Walk the dog"`
	Description  *string   `json:"description" validate:"omitempty,min=1,max=1000" example:"Morning walk at 06:30"`
	AuthorUserID uuid.UUID `json:"author_user_id" validate:"required" example:"279ee73e-0132-4d5e-9498-cfe6fb43742c"`
}

type createTaskResponse taskDTOResponse

// CreateTask godoc
// @Summary Create a task
// @Description Creates a new task in the system.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body createTaskRequest true "Task creation payload"
// @Success 201 {object} createTaskResponse "Created task"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "Author not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	var req createTaskRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		respHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	serviceParams := tasks_ports_in.NewCreateTaskParams(
		req.Title,
		req.Description,
		req.AuthorUserID,
	)
	serviceResult, err := h.tasksService.CreateTask(ctx, serviceParams)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to create task")
		return
	}

	resp := createTaskResponse(taskDTOFromDomain(serviceResult.Task))

	respHandler.JSONResponse(resp, http.StatusCreated)
}
