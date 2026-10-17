package services

import "errors"

var (
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrProfileNotFound    = errors.New("profile not found")
	ErrFailedToGetProfile = errors.New("failed to get profile")
)
