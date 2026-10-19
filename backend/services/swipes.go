package services

import (
	"context"
	"errors"
	"math"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SwipesService struct {
	swipesRepository *repositories.SwipesRepository
}

func NewSwipesService(swipesRepository *repositories.SwipesRepository) *SwipesService {
	return &SwipesService{swipesRepository: swipesRepository}
}

func (s *SwipesService) CreateSwipe(ctx context.Context, swiperID uint, swipe *models.SwipePostRequest) (bool, error) {
	if err := validateSwiperID(swiperID); err != nil {
		return false, err
	}
	if err := validateSwipeTargetID(swipe.TargetUserID); err != nil {
		return false, err
	}
	if err := validateSwipeParticipants(swiperID, swipe.TargetUserID); err != nil {
		return false, err
	}
	if err := validateSwipeAction(swipe.Action); err != nil {
		return false, err
	}

	matched, err := s.swipesRepository.CreateSwipe(ctx, swiperID, swipe)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrSwipeTargetNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && pgErr.ConstraintName == "swipes_swiper_user_id_target_user_id_key":
				return false, ErrSwipeExists
			case pgErr.Code == "23503" && pgErr.ConstraintName == "swipes_swiper_user_id_fkey":
				return false, ErrInvalidSwiperID
			case pgErr.Code == "23503" && pgErr.ConstraintName == "swipes_target_user_id_fkey":
				return false, ErrSwipeTargetNotFound
			}
		}
		return false, ErrFailedToCreateSwipe
	}
	return matched, nil
}

func validateSwiperID(swiperID uint) error {
	if swiperID == 0 || swiperID > math.MaxInt32 {
		return ErrInvalidSwiperID
	}
	return nil
}

func validateSwipeTargetID(targetID uint) error {
	if targetID == 0 || targetID > math.MaxInt32 {
		return ErrInvalidSwipeTargetID
	}
	return nil
}

func validateSwipeParticipants(swiperID, targetID uint) error {
	if swiperID == targetID {
		return ErrSwipeSelf
	}
	return nil
}

func validateSwipeAction(action string) error {
	if action != "patinka" && action != "nepatinka" {
		return ErrSwipeActionInvalid
	}
	return nil
}
