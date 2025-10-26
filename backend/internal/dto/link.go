package dto

import "time"

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	Name        string `json:"name" validate:"required"`
}

type CreateLinkResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type LinkInfo struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	Transitions int64     `json:"transitions_count"`
}

type GetAllLinksResponse struct {
	Links   []LinkInfo `json:"links"`
	Message string     `json:"message"`
}

type GetLinkResponse struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	Hash        string    `json:"hash"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Background  string    `json:"background"`
	Smoothing   *float64  `json:"smoothing"`
}

type EditLinkRequest struct {
	OriginalURL string  `json:"original_url" validate:"required,url"`
	Color       string  `json:"color" validate:"required,hexadecimal,len=6"`
	Background  string  `json:"background" validate:"required,hexadecimal,len=6"`
	Smoothing   float64 `json:"smoothing" validate:"gte=0,lte=0.5"`
}

type EditLinkResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}
