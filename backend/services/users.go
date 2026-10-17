package services

import (
	"context"
	"errors"
	"math"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{usersRepository: usersRepository}
}

func (s *UsersService) GetUserProfile(ctx context.Context, userID uint) (models.UserProfile, error) {
	if err := validateUserID(userID); err != nil {
		return models.UserProfile{}, ErrProfileNotFound
	}

	profile, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserProfile{}, ErrProfileNotFound
		}
		return models.UserProfile{}, ErrFailedToGetProfile
	}

	return profile, nil
}

func validateUserID(userID uint) error {
	if userID == 0 || userID > math.MaxInt32 {
		return ErrInvalidUserID
	}
	return nil
}
