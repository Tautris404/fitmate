package repositories

import (
	"context"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewsRepository struct {
	db *pgxpool.Pool
}

func NewReviewsRepository(db *pgxpool.Pool) *ReviewsRepository {
	return &ReviewsRepository{db: db}
}

func (r *ReviewsRepository) GetReviewByID(ctx context.Context, reviewID uint) (models.Review, error) {
	var review models.Review

	err := r.db.QueryRow(ctx, `
		SELECT r.id, u.first_name, COALESCE(r.comment, ''), r.rating
		FROM reviews r
		JOIN users u ON u.id = r.author_id
		WHERE r.id = $1`, reviewID,
	).Scan(&review.ID, &review.AuthorName, &review.Comment, &review.Rating)
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (r *ReviewsRepository) CreateReview(ctx context.Context, authorID uint, review *models.ReviewPostRequest) (uint, error) {
	var reviewID uint

	err := r.db.QueryRow(ctx, `
		INSERT INTO reviews (trainer_user_id, author_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		review.TrainerUserID, authorID, review.Rating, review.Comment,
	).Scan(&reviewID)
	if err != nil {
		return 0, err
	}

	return reviewID, nil
}

func (r *ReviewsRepository) UpdateReview(ctx context.Context, reviewID, authorID uint, review *models.ReviewUpdateRequest) error {
	result, err := r.db.Exec(ctx, `
		UPDATE reviews
		SET rating = $1, comment = $2
		WHERE id = $3 AND author_id = $4`,
		review.Rating, review.Comment, reviewID, authorID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ReviewsRepository) DeleteReview(ctx context.Context, reviewID, authorID uint) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM reviews
		WHERE id = $1 AND author_id = $2`, reviewID, authorID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
