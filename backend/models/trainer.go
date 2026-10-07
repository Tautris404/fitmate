package models

type TrainerProfileRequest struct {
	WorkExperienceYears uint   `json:"work_experience_years"`
	Education           string `json:"education"`
	SpecializationIDs   []uint `json:"specialization_ids"`
}
