package repositories

import (
	"context"

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

func (r *GymsRepository) GetGymByID(ctx context.Context, gymID uint) (models.Gym, error) {
	var gym models.Gym

	err := r.db.QueryRow(ctx, `
		SELECT g.id, g.name, g.address, g.city
		FROM gyms g
		WHERE g.id = $1`, gymID,
	).Scan(&gym.ID, &gym.Name, &gym.Address, &gym.City)
	if err != nil {
		return models.Gym{}, err
	}

	return gym, nil
}

func (r *GymsRepository) CreateGym(ctx context.Context, gym *models.GymRequest) (uint, error) {
	var gymID uint

	err := r.db.QueryRow(ctx, `
		INSERT INTO gyms (name, address, city)
		VALUES ($1, $2, $3)
		RETURNING id`,
		gym.Name, gym.Address, gym.City,
	).Scan(&gymID)
	if err != nil {
		return 0, err
	}

	return gymID, nil
}

func (r *GymsRepository) UpdateGym(ctx context.Context, gymID uint, gym *models.GymRequest) error {
	result, err := r.db.Exec(ctx, `
		UPDATE gyms
		SET name = $1, address = $2, city = $3
		WHERE id = $4`,
		gym.Name, gym.Address, gym.City, gymID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *GymsRepository) DeleteGym(ctx context.Context, gymID uint) error {
	result, err := r.db.Exec(ctx, `DELETE FROM gyms WHERE id = $1`, gymID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
