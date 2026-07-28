package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrEmailTaken         = errors.New("email taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCycleDetected      = errors.New("cycle detected")
	ErrInviteNotFound     = errors.New("invite not found")
	ErrInviteExpired      = errors.New("invite expired")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrPasswordTooShort   = errors.New("password too short")
)
