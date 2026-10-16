package repositories

import (
	"context"
	"fmt"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUserByID(ctx context.Context, userID uint) (models.UserProfile, error) {
	var user models.UserProfile

	err := r.db.QueryRow(ctx, `
		SELECT email, role_id, first_name, gender, COALESCE(bio, '')
		FROM users
		WHERE id = $1`,
		userID,
	).Scan(&user.Email, &user.RoleID, &user.FirstName, &user.Gender, &user.Bio)
	if err != nil {
		return models.UserProfile{}, err
	}

	switch user.RoleID {
	case 1:
		profile, err := r.GetAthleteProfile(ctx, userID)
		if err != nil {
			return models.UserProfile{}, fmt.Errorf("get athlete profile: %w", err)
		}
		user.AthleteProfile = &profile
	case 2:
		profile, err := r.GetTrainerProfile(ctx, userID)
		if err != nil {
			return models.UserProfile{}, fmt.Errorf("get trainer profile: %w", err)
		}
		user.TrainerProfile = &profile
	}

	user.Gyms, err = r.GetUserGyms(ctx, userID)
	if err != nil {
		return models.UserProfile{}, fmt.Errorf("get user gyms: %w", err)
	}
	user.PreferredWorkoutTimes, err = r.GetPreferredWorkoutTimes(ctx, userID)
	if err != nil {
		return models.UserProfile{}, fmt.Errorf("get preferred workout times: %w", err)
	}

	return user, nil
}

func (r *UsersRepository) GetAthleteProfile(ctx context.Context, userID uint) (models.AthleteProfile, error) {
	var profile models.AthleteProfile

	err := r.db.QueryRow(ctx, `
		SELECT TO_CHAR(date_of_birth, 'YYYY-MM-DD'), height_cm, weight_kg, experience_years
		FROM athlete_profiles
		WHERE user_id = $1`,
		userID,
	).Scan(&profile.DateOfBirth, &profile.Height, &profile.Weight, &profile.ExperienceYears)
	if err != nil {
		return models.AthleteProfile{}, err
	}
	return profile, nil
}

func (r *UsersRepository) GetTrainerProfile(ctx context.Context, userID uint) (models.TrainerProfileDetailed, error) {
	var profile models.TrainerProfileDetailed
	var trainerID int64

	err := r.db.QueryRow(ctx, `
		SELECT id, work_experience_years, education
		FROM trainer_profiles
		WHERE user_id = $1`,
		userID,
	).Scan(&trainerID, &profile.WorkExperienceYears, &profile.Education)
	if err != nil {
		return models.TrainerProfileDetailed{}, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.name
		FROM trainer_specializations ts
		JOIN specializations s ON s.id = ts.specialization_id
		WHERE ts.trainer_id = $1
		ORDER BY s.id`,
		trainerID,
	)
	if err != nil {
		return models.TrainerProfileDetailed{}, err
	}
	defer rows.Close()

	profile.Specializations = make([]models.Specialization, 0)
	for rows.Next() {
		var specialization models.Specialization
		if err := rows.Scan(&specialization.ID, &specialization.Name); err != nil {
			return models.TrainerProfileDetailed{}, err
		}
		profile.Specializations = append(profile.Specializations, specialization)
	}
	if err := rows.Err(); err != nil {
		return models.TrainerProfileDetailed{}, err
	}

	return profile, nil
}

func (r *UsersRepository) GetUserGyms(ctx context.Context, userID uint) ([]models.Gym, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id, g.name, g.address, g.city
		FROM user_gyms ug
		JOIN gyms g ON g.id = ug.gym_id
		WHERE ug.user_id = $1
		ORDER BY g.id`, userID)
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

func (r *UsersRepository) GetPreferredWorkoutTimes(ctx context.Context, userID uint) ([]models.PreferredWorkoutTime, error) {
	rows, err := r.db.Query(ctx, `
		SELECT day_of_week, TO_CHAR(start_time, 'HH24:MI'), TO_CHAR(end_time, 'HH24:MI')
		FROM user_preferred_workout_times
		WHERE user_id = $1
		ORDER BY day_of_week, start_time, end_time`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workoutTimes := make([]models.PreferredWorkoutTime, 0)
	for rows.Next() {
		var workoutTime models.PreferredWorkoutTime
		if err := rows.Scan(&workoutTime.DayOfWeek, &workoutTime.StartTime, &workoutTime.EndTime); err != nil {
			return nil, err
		}
		workoutTimes = append(workoutTimes, workoutTime)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workoutTimes, nil
}
