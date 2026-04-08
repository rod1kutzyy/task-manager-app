package users_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getUserResponse userDTOResponse

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	userDomain, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get user")
		return
	}

	resp := getUserResponse(userDTOFromDomain(userDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}
