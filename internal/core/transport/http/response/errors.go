package response

type ErrorResponse struct {
	Error   string `json:"error" example:"internal server error"`
	Message string `json:"message" example:"human-readable message"`
}
