package service

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource already exists")
	ErrUnauthorized   = errors.New("invalid credentials")
	ErrForbidden      = errors.New("access denied")
	ErrInvalidToken   = errors.New("invalid or expired token")
	ErrInternalServer = errors.New("internal server error")
)
