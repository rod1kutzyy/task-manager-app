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
	FullName    http_types.Nullable[string] `json:"full_name"`
	PhoneNumber http_types.Nullable[string] `json:"phone_number"`
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

func (h *handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	userID, err := request.GetIntPathValue(r, "id")
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
