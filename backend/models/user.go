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
	Email                 string                  `json:"email"`
	RoleID                uint                    `json:"role_id"`
	FirstName             string                  `json:"first_name"`
	Gender                string                  `json:"gender"`
	Bio                   string                  `json:"bio"`
	AthleteProfile        *AthleteProfile         `json:"athlete_profile,omitempty"`
	TrainerProfile        *TrainerProfileDetailed `json:"trainer_profile,omitempty"`
	Gyms                  []Gym                   `json:"gyms"`
	PreferredWorkoutTimes []PreferredWorkoutTime  `json:"preferred_workout_times"`
}

type Gym struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}
