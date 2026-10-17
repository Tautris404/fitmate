package services

import (
	"context"
	"errors"
	"math"
	"strings"
	"unicode/utf8"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5"
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

func (s *GymsService) GetGymByID(ctx context.Context, gymID uint) (models.Gym, error) {
	if err := validateGymID(gymID); err != nil {
		return models.Gym{}, ErrGymNotFound
	}

	gym, err := s.gymsRepository.GetGymByID(ctx, gymID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Gym{}, ErrGymNotFound
		}
		return models.Gym{}, ErrFailedToGetGym
	}

	return gym, nil
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

func (s *GymsService) UpdateGym(ctx context.Context, gymID uint, gym *models.GymRequest) error {
	if err := validateGymID(gymID); err != nil {
		return ErrGymNotFound
	}

	gym.Name = strings.TrimSpace(gym.Name)
	gym.Address = strings.TrimSpace(gym.Address)
	gym.City = strings.TrimSpace(gym.City)

	if err := validateGymName(gym.Name); err != nil {
		return err
	}
	if err := validateGymAddress(gym.Address); err != nil {
		return err
	}
	if err := validateGymCity(gym.City); err != nil {
		return err
	}

	if err := s.gymsRepository.UpdateGym(ctx, gymID, gym); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrGymNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "gyms_name_address_city_key" {
			return ErrGymExists
		}

		return ErrFailedToUpdateGym
	}

	return nil
}

func (s *GymsService) DeleteGym(ctx context.Context, gymID uint) error {
	if err := validateGymID(gymID); err != nil {
		return ErrGymNotFound
	}

	if err := s.gymsRepository.DeleteGym(ctx, gymID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrGymNotFound
		}
		return ErrFailedToDeleteGym
	}

	return nil
}

func validateGymID(gymID uint) error {
	if gymID == 0 || gymID > math.MaxInt32 {
		return ErrInvalidGymID
	}
	return nil
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
