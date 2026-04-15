package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
	http_types "github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/types"
)

type patchTaskRequest struct {
	Title       http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Walk the dog"`
	Description http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Morning walk at 06:30"`
	Completed   http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

func (r *patchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can not be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can not be null")
		}
	}

	return nil
}

type patchTaskResponse taskDTOResponse

// PatchTask godoc
// @Summary Partially update a task
// @Description Updates task fields using three-state semantics for each field.
// @Description 1. Field is omitted: the value is not changed.
// @Description 2. Field has a value: the value is updated.
// @Description 3. Field is explicitly null: the value is cleared (set to NULL).
// @Description Constraints: `title` and `completed` cannot be set to null.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body patchTaskRequest true "Task patch payload"
// @Success 200 {object} patchTaskResponse "Updated task"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *handler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	taskID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	var req patchTaskRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		respHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(req)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	resp := patchTaskResponse(taskDTOFromDomain(taskDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}

func taskPatchFromRequest(req patchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		req.Title.ToDomain(),
		req.Description.ToDomain(),
		req.Completed.ToDomain(),
	)
}
