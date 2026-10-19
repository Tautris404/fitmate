package services

import "errors"

var (
	ErrFailedToGetAthletes   = errors.New("failed to get athletes")
	ErrInvalidAthleteSeed    = errors.New("invalid athlete seed")
	ErrInvalidAthleteAfterID = errors.New("invalid athlete after ID")
)
