package services

import "errors"

var (
	ErrProfileNotFound    = errors.New("profile not found")
	ErrFailedToGetProfile = errors.New("failed to get profile")
)
