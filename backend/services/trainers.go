package services

import (
	"context"

	"fitmate/backend/models"
	"fitmate/backend/repositories"
)

type TrainersService struct {
	trainersRepository *repositories.TrainersRepository
}

func NewTrainersService(trainersRepository *repositories.TrainersRepository) *TrainersService {
	return &TrainersService{trainersRepository: trainersRepository}
}

func (s *TrainersService) GetTrainerReviews(ctx context.Context, trainerUserID uint) ([]models.Review, error) {
	if err := validateReviewTrainerID(trainerUserID); err != nil {
		return nil, err
	}

	reviews, err := s.trainersRepository.GetTrainerReviews(ctx, trainerUserID)
	if err != nil {
		return nil, ErrFailedToGetTrainerReviews
	}

	return reviews, nil
}
