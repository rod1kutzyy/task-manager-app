package users_transport_http

import (
	"net/http"
)

type createUserRequest struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}
type createUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
