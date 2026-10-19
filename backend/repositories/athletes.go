package repositories

import (
	"context"
	"errors"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AthletesRepository struct {
	db *pgxpool.Pool
}

func NewAthletesRepository(db *pgxpool.Pool) *AthletesRepository {
	return &AthletesRepository{db: db}
}

func (r *AthletesRepository) GetAthletes(ctx context.Context, userID uint, filters models.AthleteFilters, seed string, afterID uint, limit int) ([]models.AthleteListItem, bool, error) {
	if limit < 1 || limit > 100 {
		return nil, false, errors.New("limit must be between 1 and 100")
	}

	rows, err := r.db.Query(ctx, `
		WITH athlete_data AS (
			SELECT
				users.id,
				users.first_name,
				users.gender,
				COALESCE(users.bio, '') AS bio,
				EXTRACT(YEAR FROM AGE(CURRENT_DATE, profile.date_of_birth))::int AS age,
				profile.height_cm,
				profile.weight_kg,
				profile.experience_years,
				md5(users.id::text || ':' || @seed::text) COLLATE "C" AS random_order
			FROM users
			JOIN athlete_profiles profile ON profile.user_id = users.id
			WHERE users.role_id = 1
				AND users.id <> @user_id
		), matching_athletes AS (
			SELECT *
			FROM athlete_data athlete
			WHERE athlete.age BETWEEN @min_age AND @max_age
				AND athlete.experience_years BETWEEN @min_experience AND @max_experience
				AND (@preferred_gender::text = 'visi' OR athlete.gender = @preferred_gender)
				AND NOT EXISTS (
					SELECT 1
					FROM swipes
					WHERE swiper_user_id = @user_id
						AND target_user_id = athlete.id
				)
		), selected_athletes AS (
			SELECT *
			FROM matching_athletes
			WHERE @after_id::int = 0
				OR (random_order, id) > (
					md5(@after_id::text || ':' || @seed::text) COLLATE "C",
					@after_id::int
				)
			ORDER BY random_order, id
			LIMIT @fetch_limit
		)
		SELECT
			athlete.id,
			athlete.first_name,
			athlete.bio,
			athlete.age,
			athlete.height_cm,
			athlete.weight_kg,
			athlete.experience_years,
			COALESCE((
				SELECT jsonb_agg(
					jsonb_build_object(
						'id', image.id,
						'image_url', image.image_url
					) ORDER BY image.id
				)
				FROM user_images image
				WHERE image.user_id = athlete.id
			), '[]'::jsonb) AS images,
			COALESCE((
				SELECT jsonb_agg(
					jsonb_build_object(
						'id', gym.id,
						'name', gym.name,
						'address', gym.address,
						'city', gym.city
					) ORDER BY gym.id
				)
				FROM user_gyms
				JOIN gyms gym ON gym.id = user_gyms.gym_id
				WHERE user_gyms.user_id = athlete.id
			), '[]'::jsonb) AS gyms,
			COALESCE((
				SELECT jsonb_agg(
					jsonb_build_object(
						'day_of_week', workout.day_of_week,
						'start_time', TO_CHAR(workout.start_time, 'HH24:MI'),
						'end_time', TO_CHAR(workout.end_time, 'HH24:MI')
					) ORDER BY workout.day_of_week, workout.start_time, workout.end_time
				)
				FROM user_preferred_workout_times workout
				WHERE workout.user_id = athlete.id
			), '[]'::jsonb) AS preferred_workout_times
		FROM selected_athletes athlete
		ORDER BY athlete.random_order, athlete.id`, pgx.NamedArgs{
		"user_id":          userID,
		"seed":             seed,
		"min_age":          filters.MinAge,
		"max_age":          filters.MaxAge,
		"preferred_gender": filters.PreferredGender,
		"min_experience":   filters.MinExperience,
		"max_experience":   filters.MaxExperience,
		"after_id":         afterID,
		"fetch_limit":      limit + 1,
	})
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	athletes := make([]models.AthleteListItem, 0)
	for rows.Next() {
		var athlete models.AthleteListItem
		if err := rows.Scan(&athlete.ID, &athlete.FirstName, &athlete.Bio,
			&athlete.Age, &athlete.Height, &athlete.Weight, &athlete.ExperienceYears,
			&athlete.Images, &athlete.Gyms, &athlete.PreferredWorkoutTimes); err != nil {
			return nil, false, err
		}
		athletes = append(athletes, athlete)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(athletes) > limit
	if hasMore {
		athletes = athletes[:limit]
	}
	return athletes, hasMore, nil
}
