package models

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RoleID    uint   `json:"role"`
	FirstName string `json:"first_name"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`

	GymIDs                []uint                 `json:"gym_ids"`
	PreferredWorkoutTimes []PreferredWorkoutTime `json:"preferred_workout_times"`

	AthleteProfile *AthleteProfileRequest `json:"athlete_profile,omitempty"`
	TrainerProfile *TrainerProfileRequest `json:"trainer_profile,omitempty"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}