package repositories

import (
	"context"
	"errors"
	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserCommitUncertain = errors.New("user commit outcome is uncertain")

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GymExists(ctx context.Context, gymID uint) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM gyms WHERE id = $1)", gymID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepository) SpecializationExists(ctx context.Context, specializationID uint) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM specializations WHERE id = $1)", specializationID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.RegisterRequest, passwordHash string, imageURLs []string) (uint, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userID uint

	err = tx.QueryRow(ctx, `INSERT INTO users (email, password_hash, role_id, first_name, gender, bio)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
	user.Email, passwordHash, user.RoleID, user.FirstName, user.Gender, user.Bio).Scan(&userID)

	if err != nil {
		return 0, err
	}

	if user.RoleID == 1 && user.AthleteProfile != nil {
		err = createAthleteProfile(ctx, tx, userID, user.AthleteProfile)
		if err != nil {
			return 0, err
		}
	} else if user.RoleID == 2 && user.TrainerProfile != nil {
		err = createTrainerProfile(ctx, tx, userID, user.TrainerProfile)
		if err != nil {
			return 0, err
		}
	}

	if err := addUserPreferredWorkoutTimes(ctx, tx, userID, user.PreferredWorkoutTimes); err != nil {
		return 0, err
	}
	if err := addUserGyms(ctx, tx, userID, user.GymIDs); err != nil {
		return 0, err
	}
	if err := addUserImages(ctx, tx, userID, imageURLs); err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return 0, err
		}
		return 0, errors.Join(ErrUserCommitUncertain, err)
	}
	return userID, nil
}

func createAthleteProfile(ctx context.Context, tx pgx.Tx, userID uint, profile *models.AthleteProfileRequest) error {
	_, err := tx.Exec(ctx, `INSERT INTO athlete_profiles (user_id, date_of_birth, height_cm, weight_kg, experience_years)
	VALUES ($1, $2::date, $3, $4, $5)`,
	userID, profile.DateOfBirth, profile.Height, profile.Weight, profile.ExperienceYears)

	if err != nil {
		return err
	}
	return nil
}

func createTrainerProfile(ctx context.Context, tx pgx.Tx, userID uint, profile *models.TrainerProfileRequest) error {
	var trainerID uint

	err := tx.QueryRow(ctx, `INSERT INTO trainer_profiles (user_id, work_experience_years, education)
	VALUES ($1, $2, $3)
	RETURNING id`,
	userID, profile.WorkExperienceYears, profile.Education).Scan(&trainerID)

	if err != nil {
		return err
	}

	err = addTrainerSpecializations(ctx, tx, trainerID, profile.SpecializationIDs)
	if err != nil {
		return err
	}

	return nil
}

func addTrainerSpecializations(ctx context.Context, tx pgx.Tx, trainerID uint, specializationIDs []uint) error {
	for _, specializationID := range specializationIDs {
		_, err := tx.Exec(ctx, `INSERT INTO trainer_specializations (trainer_id, specialization_id)
		VALUES ($1, $2)`, trainerID, specializationID)

		if err != nil {
			return err
		}
	}
	return nil
}

func addUserPreferredWorkoutTimes(ctx context.Context, tx pgx.Tx, userID uint, preferredTimes []models.PreferredWorkoutTime) error {
	for _, time := range preferredTimes {
		_, err := tx.Exec(ctx, `INSERT INTO user_preferred_workout_times (user_id, day_of_week, start_time, end_time)
		VALUES ($1, $2, $3::time, $4::time)`, userID, time.DayOfWeek, time.StartTime, time.EndTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func addUserGyms(ctx context.Context, tx pgx.Tx, userID uint, gymIDs []uint) error {
	for _, gymID := range gymIDs {
		_, err := tx.Exec(ctx, `INSERT INTO user_gyms (user_id, gym_id) VALUES ($1, $2)`, userID, gymID)
		if err != nil {
			return err
		}
	}
	return nil
}

func addUserImages(ctx context.Context, tx pgx.Tx, userID uint, imageURLs []string) error {
	for _, imageURL := range imageURLs {
		_, err := tx.Exec(ctx, `INSERT INTO user_images (user_id, image_url) VALUES ($1, $2)`, userID, imageURL)
		if err != nil {
			return err
		}
	}
	return nil
}
