package dtos

// ErrorResponse represents a standard error object returned by the API.
type ErrorResponse struct {
	Message string `json:"message" example:"error description"`
	Code    int    `json:"code" example:"400"`
}
