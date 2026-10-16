package repositories

import (
	"context"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID uint,
) (models.UserProfile, error) {
	var user models.UserProfile

	err := r.db.QueryRow(ctx, `
		SELECT id, email, role_id, first_name, gender, COALESCE(bio, '')
		FROM users
		WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.RoleID,
		&user.FirstName,
		&user.Gender,
		&user.Bio,
	)
	if err != nil {
		return models.UserProfile{}, err
	}

	return user, nil
}
