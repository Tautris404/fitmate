package repositories

import (
	"context"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SwipesRepository struct {
	db *pgxpool.Pool
}

func NewSwipesRepository(db *pgxpool.Pool) *SwipesRepository {
	return &SwipesRepository{db: db}
}

func (r *SwipesRepository) CreateSwipe(ctx context.Context, swiperID uint, swipe *models.SwipePostRequest) (bool, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	if err := lockSwipePair(ctx, tx, swiperID, swipe.TargetUserID); err != nil {
		return false, err
	}
	if err := insertSwipe(ctx, tx, swiperID, swipe); err != nil {
		return false, err
	}
	var matched bool
	if swipe.Action == "patinka" {
		matched, err = createChatForMutualLike(ctx, tx, swiperID, swipe.TargetUserID)
		if err != nil {
			return false, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return matched, nil
}

func lockSwipePair(ctx context.Context, tx pgx.Tx, swiperID, targetID uint) error {
	_, err := tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(
			LEAST($1::int, $2::int),
			GREATEST($1::int, $2::int)
		)`, swiperID, targetID)
	return err
}

func insertSwipe(ctx context.Context, tx pgx.Tx, swiperID uint, swipe *models.SwipePostRequest) error {
	result, err := tx.Exec(ctx, `
  INSERT INTO swipes (swiper_user_id, target_user_id, action)
  SELECT $1, profile.user_id, $3
  FROM athlete_profiles profile
  JOIN users ON users.id = profile.user_id AND users.role_id = 1
  WHERE profile.user_id = $2`, swiperID, swipe.TargetUserID, swipe.Action)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func createChatForMutualLike(ctx context.Context, tx pgx.Tx, swiperID, targetID uint) (bool, error) {
	liked, err := hasReverseLike(ctx, tx, swiperID, targetID)
	if err != nil {
		return false, err
	}
	if !liked {
		return false, nil
	}
	if err := createMatchChat(ctx, tx, swiperID, targetID); err != nil {
		return false, err
	}
	return true, nil
}

func hasReverseLike(ctx context.Context, tx pgx.Tx, swiperID, targetID uint) (bool, error) {
	var liked bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM swipes
			WHERE swiper_user_id = $1 AND target_user_id = $2 AND action = 'patinka'
		)`, targetID, swiperID).Scan(&liked)
	return liked, err
}

func createMatchChat(ctx context.Context, tx pgx.Tx, userID, otherUserID uint) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO chats (user1_id, user2_id, type)
		VALUES (LEAST($1::int, $2::int), GREATEST($1::int, $2::int), 'match')
		ON CONFLICT (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id)) DO NOTHING`, userID, otherUserID)
	return err
}
