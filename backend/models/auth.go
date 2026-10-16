package models

import "time"

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RoleID    uint   `json:"role"`
	FirstName string `json:"first_name"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`

	GymIDs                []uint                 `json:"gym_ids"`
	PreferredWorkoutTimes []PreferredWorkoutTime `json:"preferred_workout_times"`

	AthleteProfile *AthleteProfile `json:"athlete_profile,omitempty"`
	TrainerProfile *TrainerProfile `json:"trainer_profile,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}
