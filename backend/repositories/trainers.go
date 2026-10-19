package repositories

import (
	"context"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TrainersRepository struct {
	db *pgxpool.Pool
}

func NewTrainersRepository(db *pgxpool.Pool) *TrainersRepository {
	return &TrainersRepository{db: db}
}

func (r *TrainersRepository) GetTrainerReviews(ctx context.Context, trainerUserID uint) ([]models.Review, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, u.first_name, COALESCE(r.comment, ''), r.rating
		FROM reviews r
		JOIN users u ON u.id = r.author_id
		WHERE r.trainer_user_id = $1
		ORDER BY r.created_at DESC, r.id DESC`, trainerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]models.Review, 0)
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(
			&review.ID,
			&review.AuthorName,
			&review.Comment,
			&review.Rating,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
