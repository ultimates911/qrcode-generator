package dto

type GenerateQRCodeRequest struct {
	URL        string  `json:"url" validate:"required,url"`
	Color      string  `json:"color" validate:"required"`
	Background string  `json:"background" validate:"required"`
	Smoothing  float64 `json:"smoothing" validate:"gte=0,lte=0.5"`
}
