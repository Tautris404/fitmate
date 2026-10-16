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
type AuthUser struct {
	ID           uint
	PasswordHash string
	RoleID       uint
}
type UserProfile struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	RoleID    uint   `json:"role_id"`
	FirstName string `json:"first_name"`
	Gender    string `json:"gender"`
	Bio       string `json:"bio"`
}
