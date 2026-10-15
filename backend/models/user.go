package models

import (
	"time"
)

type PreferredWorkoutTime struct {
	DayOfWeek uint   `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ParsedPreferredWorkoutTime struct {
	DayOfWeek uint
	StartTime time.Time
	EndTime   time.Time
}