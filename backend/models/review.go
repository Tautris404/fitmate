package models

type Review struct {
	ID         uint   `json:"id"`
	AuthorName string `json:"author_name"`
	Comment    string `json:"comment"`
	Rating     uint   `json:"rating"`
}

type ReviewPostRequest struct {
	TrainerUserID uint   `json:"trainer_user_id"`
	Comment       string `json:"comment"`
	Rating        uint   `json:"rating"`
}

type ReviewUpdateRequest struct {
	Comment string `json:"comment"`
	Rating  uint   `json:"rating"`
}

type ReviewResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
