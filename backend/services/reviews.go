package services

import (
	"context"
	"errors"
	"math"
	"strings"

	"fitmate/backend/models"
	"fitmate/backend/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ReviewsService struct {
	reviewsRepository *repositories.ReviewsRepository
}

func NewReviewsService(reviewsRepository *repositories.ReviewsRepository) *ReviewsService {
	return &ReviewsService{reviewsRepository: reviewsRepository}
}

func (s *ReviewsService) GetReviewByID(ctx context.Context, reviewID uint) (models.Review, error) {
	if err := validateReviewID(reviewID); err != nil {
		return models.Review{}, ErrReviewNotFound
	}

	review, err := s.reviewsRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Review{}, ErrReviewNotFound
		}
		return models.Review{}, ErrFailedToGetReview
	}

	return review, nil
}

func (s *ReviewsService) CreateReview(ctx context.Context, authorID uint, review *models.ReviewPostRequest) (uint, error) {
	if err := validateReviewAuthorID(authorID); err != nil {
		return 0, err
	}
	if err := validateReviewTrainerID(review.TrainerUserID); err != nil {
		return 0, err
	}
	if err := validateReviewRating(review.Rating); err != nil {
		return 0, err
	}

	review.Comment = strings.TrimSpace(review.Comment)

	reviewID, err := s.reviewsRepository.CreateReview(ctx, authorID, review)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && pgErr.ConstraintName == "reviews_trainer_user_id_author_id_key":
				return 0, ErrReviewExists
			case pgErr.Code == "23503" && pgErr.ConstraintName == "reviews_trainer_user_id_fkey":
				return 0, ErrTrainerNotFound
			case pgErr.Code == "23503" && pgErr.ConstraintName == "reviews_author_id_fkey":
				return 0, ErrInvalidReviewAuthorID
			}
		}

		return 0, ErrFailedToCreateReview
	}

	return reviewID, nil
}

func (s *ReviewsService) UpdateReview(ctx context.Context, reviewID, authorID uint, review *models.ReviewUpdateRequest) error {
	if err := validateReviewID(reviewID); err != nil {
		return err
	}
	if err := validateReviewAuthorID(authorID); err != nil {
		return err
	}
	if err := validateReviewRating(review.Rating); err != nil {
		return err
	}

	review.Comment = strings.TrimSpace(review.Comment)

	if err := s.reviewsRepository.UpdateReview(ctx, reviewID, authorID, review); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrReviewNotFound
		}
		return ErrFailedToUpdateReview
	}

	return nil
}

func (s *ReviewsService) DeleteReview(ctx context.Context, reviewID, authorID uint) error {
	if err := validateReviewID(reviewID); err != nil {
		return err
	}
	if err := validateReviewAuthorID(authorID); err != nil {
		return err
	}

	if err := s.reviewsRepository.DeleteReview(ctx, reviewID, authorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrReviewNotFound
		}
		return ErrFailedToDeleteReview
	}

	return nil
}

func validateReviewAuthorID(authorID uint) error {
	if authorID == 0 || authorID > math.MaxInt32 {
		return ErrInvalidReviewAuthorID
	}
	return nil
}

func validateReviewTrainerID(trainerUserID uint) error {
	if trainerUserID == 0 || trainerUserID > math.MaxInt32 {
		return ErrInvalidTrainerID
	}
	return nil
}

func validateReviewID(reviewID uint) error {
	if reviewID == 0 || reviewID > math.MaxInt32 {
		return ErrInvalidReviewID
	}
	return nil
}

func validateReviewRating(rating uint) error {
	if rating < 1 || rating > 5 {
		return ErrReviewRatingInvalid
	}
	return nil
}
