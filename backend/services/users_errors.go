package services

import "errors"

var (
	ErrPasswordTooLong       = errors.New("password must not exceed 72 bytes")
	ErrInvalidPage           = errors.New("invalid page")
	ErrInvalidLimit          = errors.New("limit must be between 1 and 100")
	ErrFailedToGetUsers      = errors.New("failed to get users")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrProfileNotFound       = errors.New("profile not found")
	ErrFailedToGetProfile    = errors.New("failed to get profile")
	ErrFailedToUpdateProfile = errors.New("failed to update profile")
	ErrFailedToDeleteProfile = errors.New("failed to delete profile")
	ErrImageLimitExceeded    = errors.New("profile cannot have more than 10 images")
	ErrFailedToAddImage      = errors.New("failed to add image")
	ErrSingleImageRequired   = errors.New("exactly one image is required")
	ErrInvalidImageID        = errors.New("invalid image id")
	ErrImageNotFound         = errors.New("image not found")
	ErrFailedToDeleteImage   = errors.New("failed to delete image")
	ErrLastImage             = errors.New("profile must have at least one image")
)
