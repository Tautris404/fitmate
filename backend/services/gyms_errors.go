package services

import "errors"

var (
	ErrFailedToGetGymTrainers = errors.New("failed to get gym trainers")
	ErrFailedToGetGyms        = errors.New("failed to get gyms")
	ErrFailedToGetGym         = errors.New("failed to get gym")
	ErrInvalidGymID           = errors.New("invalid gym id")
	ErrGymNotFound            = errors.New("gym not found")
	ErrGymNameEmpty           = errors.New("gym name is required")
	ErrGymNameTooLong         = errors.New("gym name field is too long")
	ErrGymAddressEmpty        = errors.New("gym address is required")
	ErrGymAddressTooLong      = errors.New("gym address field is too long")
	ErrGymCityEmpty           = errors.New("gym city is required")
	ErrGymCityTooLong         = errors.New("gym city field is too long")
	ErrFailedToCreateGym      = errors.New("failed to create gym")
	ErrFailedToUpdateGym      = errors.New("failed to update gym")
	ErrFailedToDeleteGym      = errors.New("failed to delete gym")
	ErrGymExists              = errors.New("gym already exists")
)
