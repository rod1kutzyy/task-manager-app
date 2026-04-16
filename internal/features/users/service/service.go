package users_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

type service struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error)
}

func NewService(usersRepository UsersRepository) *service {
	return &service{
		usersRepository: usersRepository,
	}
}
