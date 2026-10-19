package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"strings"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5"
)

const athletePageSize = 20

type AthletesService struct {
	athletesRepository *repositories.AthletesRepository
	usersRepository    *repositories.UsersRepository
}

func NewAthletesService(athletesRepository *repositories.AthletesRepository, usersRepository *repositories.UsersRepository) *AthletesService {
	return &AthletesService{athletesRepository: athletesRepository, usersRepository: usersRepository}
}

func (s *AthletesService) GetAthletes(ctx context.Context, userID uint, seed string, afterID uint) (models.AthleteListResponse, error) {
	if err := validateUserID(userID); err != nil {
		return models.AthleteListResponse{}, ErrProfileNotFound
	}
	if err := validateAthleteAfterID(afterID); err != nil {
		return models.AthleteListResponse{}, err
	}
	if err := validateAthleteSeed(seed, afterID); err != nil {
		return models.AthleteListResponse{}, err
	}
	seed, err := prepareAthleteSeed(seed)
	if err != nil {
		return models.AthleteListResponse{}, err
	}

	profile, err := s.usersRepository.GetAthleteProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AthleteListResponse{}, ErrProfileNotFound
		}
		return models.AthleteListResponse{}, ErrFailedToGetAthletes
	}
	filters := models.AthleteFilters{
		MinAge:          profile.MinAge,
		MaxAge:          profile.MaxAge,
		PreferredGender: profile.PreferredGender,
		MinExperience:   profile.MinExperience,
		MaxExperience:   profile.MaxExperience,
	}
	athletes, hasMore, err := s.athletesRepository.GetAthletes(ctx, userID, filters, seed, afterID, athletePageSize)
	if err != nil {
		return models.AthleteListResponse{}, ErrFailedToGetAthletes
	}
	response := models.AthleteListResponse{Athletes: athletes, Seed: seed, HasMore: hasMore}
	if hasMore {
		response.NextAfterID = athletes[len(athletes)-1].ID
	}
	return response, nil
}

func validateAthleteAfterID(afterID uint) error {
	if afterID > math.MaxInt32 {
		return ErrInvalidAthleteAfterID
	}
	return nil
}

func validateAthleteSeed(seed string, afterID uint) error {
	if seed == "" {
		if afterID != 0 {
			return ErrInvalidAthleteSeed
		}
		return nil
	}
	if len(seed) != 32 {
		return ErrInvalidAthleteSeed
	}
	if _, err := hex.DecodeString(seed); err != nil {
		return ErrInvalidAthleteSeed
	}
	return nil
}

func prepareAthleteSeed(seed string) (string, error) {
	if seed != "" {
		return strings.ToLower(seed), nil
	}
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", ErrFailedToGetAthletes
	}
	return hex.EncodeToString(value), nil
}
