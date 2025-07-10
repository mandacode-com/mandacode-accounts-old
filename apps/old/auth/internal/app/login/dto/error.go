package logindto

import "errors"

var (
	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserNotActive is returned when the user is not active.
	ErrUserNotActive = errors.New("user is not active")
	// ErrUserNotVerified is returned when the user is not verified.
	ErrUserNotVerified = errors.New("user is not verified")
)
