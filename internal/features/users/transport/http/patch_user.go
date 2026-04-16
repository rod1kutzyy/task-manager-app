package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
	http_types "github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/types"
)

type patchUserRequest struct {
	FullName    http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Ivan Ivanov"`
	PhoneNumber http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

func (r *patchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can not be null")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must start with '+' symbol")
			}
		}
	}

	return nil
}

type patchUserResponse userDTOResponse

// PatchUser godoc
// @Summary Partially update a user
// @Description Updates user fields using three-state semantics for each field.
// @Description 1. Field is omitted: the value is not changed.
// @Description 2. Field has a value: the value is updated.
// @Description 3. Field is explicitly null: the value is cleared (set to NULL).
// @Description Constraint: `full_name` cannot be set to null.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)" Format(uuid)
// @Param request body patchUserRequest true "User patch payload"
// @Success 200 {object} patchUserResponse "Updated user"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	var req patchUserRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		respHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	resp := patchUserResponse(userDTOFromDomain(userDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}

func userPatchFromRequest(req patchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		req.FullName.ToDomain(),
		req.PhoneNumber.ToDomain(),
	)
}
