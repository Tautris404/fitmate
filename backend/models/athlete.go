package models

type AthleteProfile struct {
	DateOfBirth     string  `json:"date_of_birth"`
	Height          float64 `json:"height"`
	Weight          float64 `json:"weight"`
	ExperienceYears uint    `json:"experience_years"`
	MinAge          int     `json:"min_age"`
	MaxAge          int     `json:"max_age"`
	PreferredGender string  `json:"preferred_gender"`
	MinExperience   int     `json:"min_experience"`
	MaxExperience   int     `json:"max_experience"`
}

type AthleteFilters struct {
	MinAge          int    `json:"min_age"`
	MaxAge          int    `json:"max_age"`
	PreferredGender string `json:"preferred_gender"`
	MinExperience   int    `json:"min_experience"`
	MaxExperience   int    `json:"max_experience"`
}

type AthleteListItem struct {
	ID                    uint                   `json:"id"`
	FirstName             string                 `json:"first_name"`
	Bio                   string                 `json:"bio"`
	Age                   uint                   `json:"age"`
	Height                float64                `json:"height"`
	Weight                float64                `json:"weight"`
	ExperienceYears       uint                   `json:"experience_years"`
	Images                []UserImage            `json:"images"`
	Gyms                  []Gym                  `json:"gyms"`
	PreferredWorkoutTimes []PreferredWorkoutTime `json:"preferred_workout_times"`
}

type AthleteListResponse struct {
	Athletes    []AthleteListItem `json:"athletes"`
	Seed        string            `json:"seed"`
	NextAfterID uint              `json:"next_after_id"`
	HasMore     bool              `json:"has_more"`
}
