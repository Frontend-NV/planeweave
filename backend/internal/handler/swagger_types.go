package handler

import "time"

// ErrorResponse описывает тело ошибки API.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid_request"`
}

// InviteResponse ответ при создании приглашения.
type InviteResponse struct {
	Token     string    `json:"token" example:"a1b2c3"`
	ExpiresAt time.Time `json:"expiresAt"`
}
