package dto

// dto.GenericError represents a generic error response.
type GenericError struct {
	Error string `json:"error" example:"Some error message"`
}