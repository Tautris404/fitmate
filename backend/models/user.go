package models

type PreferredWorkoutTime struct {
	DayOfWeek uint   `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}