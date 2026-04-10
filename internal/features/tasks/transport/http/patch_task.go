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
	Title       http_types.Nullable[string] `json:"title"`
	Description http_types.Nullable[string] `json:"description"`
	Completed   http_types.Nullable[bool]   `json:"completed"`
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

func (h *handler) Patchtask(w http.ResponseWriter, r *http.Request) {
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
