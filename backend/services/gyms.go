package services

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5/pgconn"
)

type GymsService struct {
	gymsRepository *repositories.GymsRepository
}

func NewGymsService(gymsRepository *repositories.GymsRepository) *GymsService {
	return &GymsService{gymsRepository: gymsRepository}
}

func (s *GymsService) GetGyms(ctx context.Context) ([]models.Gym, error) {
	gyms, err := s.gymsRepository.GetGyms(ctx)
	if err != nil {
		return nil, ErrFailedToGetGyms
	}

	return gyms, nil
}
func (s *GymsService) CreateGym(ctx context.Context, gym *models.GymRequest) (uint, error) {
	gym.Name = strings.TrimSpace(gym.Name)
	gym.Address = strings.TrimSpace(gym.Address)
	gym.City = strings.TrimSpace(gym.City)

	if err := validateGymName(gym.Name); err != nil {
		return 0, err
	}
	if err := validateGymAddress(gym.Address); err != nil {
		return 0, err
	}
	if err := validateGymCity(gym.City); err != nil {
		return 0, err
	}

	exists, err := s.gymsRepository.GymExists(ctx, *gym)

	if err != nil {
		return 0, ErrGymFailedCheck
	}
	if exists {
		return 0, ErrGymExists
	}

	gymID, err := s.gymsRepository.CreateGym(ctx, gym)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "gyms_name_address_city_key" {
			return 0, ErrGymExists
		}

		return 0, ErrFailedToCreateGym
	}
	return gymID, nil
}

func validateGymName(gymName string) error {
	if gymName == "" {
		return ErrGymNameEmpty
	}
	if utf8.RuneCountInString(gymName) > 100 {
		return ErrGymNameTooLong
	}
	return nil
}
func validateGymAddress(gymAddress string) error {
	if gymAddress == "" {
		return ErrGymAddressEmpty
	}
	if utf8.RuneCountInString(gymAddress) > 255 {
		return ErrGymAddressTooLong
	}
	return nil
}

func validateGymCity(gymCity string) error {
	if gymCity == "" {
		return ErrGymCityEmpty
	}
	if utf8.RuneCountInString(gymCity) > 100 {
		return ErrGymCityTooLong
	}
	return nil
}
