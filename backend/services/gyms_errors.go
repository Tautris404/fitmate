package services

import "errors"

var (
	ErrFailedToGetGyms   = errors.New("failed to get gyms")
	ErrGymNameEmpty      = errors.New("gym name is required")
	ErrGymNameTooLong    = errors.New("gym name field is too long")
	ErrGymAddressEmpty   = errors.New("gym address is required")
	ErrGymAddressTooLong = errors.New("gym address field is too long")
	ErrGymCityEmpty      = errors.New("gym city is required")
	ErrGymCityTooLong    = errors.New("gym city field is too long")
	ErrFailedToCreateGym = errors.New("failed to create gym")
	ErrGymExists         = errors.New("gym already exists")
	ErrGymFailedCheck    = errors.New("failed to check gym")
)
