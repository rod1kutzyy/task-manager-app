package response

type ErrorResponse struct {
	Error   string `json:"error" example:"error text"`
	Message string `json:"message" example:"human-readable message"`
}
