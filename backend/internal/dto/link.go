package dto

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
}

type CreateLinkResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
