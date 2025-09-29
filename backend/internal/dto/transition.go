package dto

import "time"

type TransitionItem struct {
	ID        int64     `json:"id"`
	Country   *string   `json:"country,omitempty"`
	City      *string   `json:"city,omitempty"`
	Referer   *string   `json:"referer,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	Browser   *string   `json:"browser,omitempty"`
	OS        *string   `json:"os,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type GetTransitionsResponse struct {
	Transitions []TransitionItem `json:"transitions"`
}
