package models

type Gym struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}
type GymRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}
type GymResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
