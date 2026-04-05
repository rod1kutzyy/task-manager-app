package users_transport_http

import (
	"context"
	"net/http"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
)

type handler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
}

func NewHandler(usersService UsersService) *handler {
	return &handler{
		usersService: usersService,
	}
}

func (h *handler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
	}
}
