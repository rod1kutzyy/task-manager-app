package users_transport_http

import (
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getUserResponse userDTOResponse

// GetUser godoc
// @Summary Get a user by ID
// @Description Returns a single user by UUID identifier.
// @Tags users
// @Produce json
// @Param id path string true "User ID (UUID)" Format(uuid)
// @Success 200 {object} getUserResponse "User found"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, err := request.GetUUIDPathValue(r, "id")
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
