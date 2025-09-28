package dto

import "time"

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
}

type CreateLinkResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type GetLinkResponse struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	Hash        string    `json:"hash"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Color       string    `json:"color"`
	Background  string    `json:"background"`
	Smoothing   *float64  `json:"smoothing"`
}