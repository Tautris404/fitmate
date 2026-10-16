package models

type TrainerProfile struct {
	WorkExperienceYears uint   `json:"work_experience_years"`
	Education           string `json:"education"`
	SpecializationIDs   []uint `json:"specialization_ids"`
}
