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

type UserImage struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
}

type UserListItem struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	RoleName  string `json:"role_name"`
}

type UserListResponse struct {
	Users      []UserListItem `json:"users"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type UserProfileUpdateRequest struct {
	FirstName             string                 `json:"first_name"`
	Gender                string                 `json:"gender"`
	Bio                   string                 `json:"bio"`
	AthleteProfile        *AthleteProfile        `json:"athlete_profile,omitempty"`
	TrainerProfile        *TrainerProfile        `json:"trainer_profile,omitempty"`
	GymIDs                []uint                 `json:"gym_ids"`
	PreferredWorkoutTimes []PreferredWorkoutTime `json:"preferred_workout_times"`
}

type UserUpdateRequest struct {
	Password string `json:"password,omitempty"`
	UserProfileUpdateRequest
}

type UserProfile struct {
	Email                 string                  `json:"email"`
	RoleID                uint                    `json:"-"`
	FirstName             string                  `json:"first_name"`
	Gender                string                  `json:"gender"`
	Bio                   string                  `json:"bio"`
	Images                []UserImage             `json:"images"`
	AthleteProfile        *AthleteProfile         `json:"athlete_profile,omitempty"`
	TrainerProfile        *TrainerProfileDetailed `json:"trainer_profile,omitempty"`
	Gyms                  []Gym                   `json:"gyms"`
	PreferredWorkoutTimes []PreferredWorkoutTime  `json:"preferred_workout_times"`
}

type UserResponse struct {
	ID      uint   `json:"id,omitempty"`
	Message string `json:"message"`
}

type UserDetails struct {
	ID       uint   `json:"id"`
	RoleName string `json:"role_name"`
	UserProfile
}
