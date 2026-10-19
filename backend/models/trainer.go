package models

type TrainerProfile struct {
	WorkExperienceYears uint   `json:"work_experience_years"`
	Education           string `json:"education"`
	SpecializationIDs   []uint `json:"specialization_ids"`
}

type TrainerProfileDetailed struct {
	WorkExperienceYears uint             `json:"work_experience_years"`
	Education           string           `json:"education"`
	Specializations     []Specialization `json:"specializations"`
}

type Specialization struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TrainerListItem struct {
	ID                    uint                   `json:"id"`
	FirstName             string                 `json:"first_name"`
	Bio                   string                 `json:"bio"`
	WorkExperienceYears   uint                   `json:"work_experience_years"`
	Education             string                 `json:"education"`
	Specializations       []Specialization       `json:"specializations"`
	Images                []UserImage            `json:"images"`
	Gyms                  []Gym                  `json:"gyms"`
	PreferredWorkoutTimes []PreferredWorkoutTime `json:"preferred_workout_times"`
	AverageRating         float64                `json:"average_rating"`
	ReviewCount           int                    `json:"review_count"`
}

type TrainerListResponse struct {
	Trainers   []TrainerListItem `json:"trainers"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}
