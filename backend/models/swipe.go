package models

type SwipePostRequest struct {
	TargetUserID uint   `json:"target_user_id"`
	Action       string `json:"action"`
}

type SwipeResponse struct {
	Match   bool   `json:"match"`
	Message string `json:"message"`
}
