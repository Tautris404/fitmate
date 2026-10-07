package models

type AthleteProfileRequest struct {
	DateOfBirth     string  `json:"date_of_birth"`
	Height          float64 `json:"height"`
	Weight          float64 `json:"weight"`
	ExperienceYears uint    `json:"experience_years"`
}