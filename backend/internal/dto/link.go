package dto

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
}

type CreateLinkResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type LinkInfo struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
}

type GetAllLinksResponse struct {
	Links   []LinkInfo `json:"links"`
	Message string     `json:"message"`
}
