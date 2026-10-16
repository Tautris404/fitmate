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
