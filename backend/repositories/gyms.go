package repositories

import (
	"context"
	"errors"
	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GymsRepository struct {
	db *pgxpool.Pool
}

func NewGymsRepository(db *pgxpool.Pool) *GymsRepository {
	return &GymsRepository{db: db}
}

func (r *GymsRepository) GetGyms(ctx context.Context) ([]models.Gym, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id, g.name, g.address, g.city
		FROM gyms g
		ORDER BY g.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gyms := make([]models.Gym, 0)
	for rows.Next() {
		var gym models.Gym
		if err := rows.Scan(&gym.ID, &gym.Name, &gym.Address, &gym.City); err != nil {
			return nil, err
		}
		gyms = append(gyms, gym)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return gyms, nil
}

func (r *GymsRepository) CreateGym(ctx context.Context, gym *models.GymRequest) (uint, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var gymID uint

	err = tx.QueryRow(ctx, `INSERT INTO gyms (name, address, city)
	VALUES ($1, $2, $3) RETURNING id`,
		gym.Name, gym.Address, gym.City).Scan(&gymID)

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return 0, err
		}
		return 0, errors.Join(ErrUserCommitUncertain, err)
	}
	return gymID, nil
}

func (r *GymsRepository) GymExists(ctx context.Context, gym models.GymRequest) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM gyms WHERE name = $1 AND address = $2 AND city = $3)", gym.Name, gym.Address, gym.City).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
