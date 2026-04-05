package users_postgres_repository

import "github.com/rod1kutzyy/task-manager-app/internal/core/domain"

type UserModel struct {
	ID          int     `db:"id"`
	Version     int     `db:"version"`
	FullName    string  `db:"full_name"`
	PhoneNumber *string `db:"phone_number"`
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return userDomains
}
