package users_service

import (
	"context"
	"fmt"

	"github.com/rod1kutzyy/task-manager-app/internal/core/domain"
)

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	user = domain.CreateUser(user.FullName, user.PhoneNumber)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.SaveUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
