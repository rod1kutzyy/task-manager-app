package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
	"github.com/rod1kutzyy/task-manager-app/internal/core/transport/http/server"
)

type handler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	ListUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	PatchUser(ctx context.Context, id uuid.UUID, patch domain.UserPatch) (domain.User, error)
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
			Handler: h.ListUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
