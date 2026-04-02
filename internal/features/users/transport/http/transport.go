package users_transport_http

type handler struct {
	usersService UsersService
}

type UsersService interface {
}

func NewHandler(usersService UsersService) *handler {
	return &handler{
		usersService: usersService,
	}
}
