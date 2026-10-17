package services

import "errors"

var (
	ErrInvalidTrainerID          = errors.New("invalid trainer id")
	ErrFailedToGetTrainerReviews = errors.New("failed to get trainer reviews")
	ErrInvalidReviewID           = errors.New("invalid review id")
	ErrReviewNotFound            = errors.New("review not found")
	ErrFailedToGetReview         = errors.New("failed to get review")
	ErrInvalidReviewAuthorID     = errors.New("invalid review author id")
	ErrReviewRatingInvalid       = errors.New("rating must be between 1 and 5")
	ErrReviewExists              = errors.New("review already exists for this trainer")
	ErrTrainerNotFound           = errors.New("trainer not found")
	ErrFailedToCreateReview      = errors.New("failed to create review")
	ErrFailedToUpdateReview      = errors.New("failed to update review")
	ErrFailedToDeleteReview      = errors.New("failed to delete review")
)
