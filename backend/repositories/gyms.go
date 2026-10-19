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

func (r *GymsRepository) GetGymTrainers(ctx context.Context, gymID uint, limit, offset int) ([]models.TrainerListItem, int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM user_gyms membership
		JOIN users ON users.id = membership.user_id AND users.role_id = 2
		JOIN trainer_profiles profile ON profile.user_id = users.id
		WHERE membership.gym_id = $1`, gymID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		WITH selected_trainers AS (
			SELECT users.id, users.first_name, COALESCE(users.bio, '') AS bio,
				profile.work_experience_years, profile.education
			FROM user_gyms membership
			JOIN users ON users.id = membership.user_id AND users.role_id = 2
			JOIN trainer_profiles profile ON profile.user_id = users.id
			WHERE membership.gym_id = $1
			ORDER BY users.id
			LIMIT $2 OFFSET $3
		)
		SELECT trainer.id, trainer.first_name, trainer.bio,
			trainer.work_experience_years, trainer.education,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object(
					'id', specialization.id, 'name', specialization.name
				) ORDER BY specialization.id)
				FROM trainer_specializations
				JOIN specializations specialization ON specialization.id = trainer_specializations.specialization_id
				WHERE trainer_specializations.trainer_user_id = trainer.id
			), '[]'::jsonb) AS specializations,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object(
					'id', image.id, 'image_url', image.image_url
				) ORDER BY image.id)
				FROM user_images image
				WHERE image.user_id = trainer.id
			), '[]'::jsonb) AS images,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object(
					'id', gym.id, 'name', gym.name, 'address', gym.address, 'city', gym.city
				) ORDER BY gym.id)
				FROM user_gyms
				JOIN gyms gym ON gym.id = user_gyms.gym_id
				WHERE user_gyms.user_id = trainer.id
			), '[]'::jsonb) AS gyms,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object(
					'day_of_week', workout.day_of_week,
					'start_time', TO_CHAR(workout.start_time, 'HH24:MI'),
					'end_time', TO_CHAR(workout.end_time, 'HH24:MI')
				) ORDER BY workout.day_of_week, workout.start_time, workout.end_time)
				FROM user_preferred_workout_times workout
				WHERE workout.user_id = trainer.id
			), '[]'::jsonb) AS preferred_workout_times,
			COALESCE((
				SELECT ROUND(AVG(review.rating), 2)
				FROM reviews review
				WHERE review.trainer_user_id = trainer.id
			), 0) AS average_rating,
			(SELECT COUNT(*) FROM reviews review WHERE review.trainer_user_id = trainer.id) AS review_count
		FROM selected_trainers trainer
		ORDER BY trainer.id`, gymID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	trainers := make([]models.TrainerListItem, 0)
	for rows.Next() {
		var trainer models.TrainerListItem
		if err := rows.Scan(&trainer.ID, &trainer.FirstName, &trainer.Bio,
			&trainer.WorkExperienceYears, &trainer.Education, &trainer.Specializations,
			&trainer.Images, &trainer.Gyms, &trainer.PreferredWorkoutTimes,
			&trainer.AverageRating, &trainer.ReviewCount); err != nil {
			return nil, 0, err
		}
		trainers = append(trainers, trainer)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return trainers, total, nil
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
