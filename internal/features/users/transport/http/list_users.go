package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/rod1kutzyy/task-manager-app/internal/core/logger"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/request"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/response"
)

type getUsersResponse []userDTOResponse

// ListUsers godoc
// @Summary List users
// @Description Returns users with optional pagination via limit and offset.
// @Tags users
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {array} userDTOResponse "Users list"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	respHandler := response.NewHTTPResponseHandler(w, logger)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	usersDomain, err := h.usersService.ListUsers(ctx, limit, offset)
	if err != nil {
		respHandler.ErrorResponse(err, "failed to list users")
		return
	}

	resp := getUsersResponse(userDTOFromDomains(usersDomain))

	respHandler.JSONResponse(resp, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
