package repositories

import (
	"context"
	"errors"

	"fitmate/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

var ErrUserImageLimitExceeded = errors.New("user image limit exceeded")
var ErrLastUserImage = errors.New("cannot delete last user image")

func NewUsersRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUsers(ctx context.Context, limit, offset int) ([]models.UserListItem, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.email, u.first_name, r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY u.id
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]models.UserListItem, 0)
	for rows.Next() {
		var user models.UserListItem
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.RoleName); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UsersRepository) DeleteUser(ctx context.Context, userID uint) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := lockUserRow(ctx, tx, userID); err != nil {
		return nil, err
	}
	imageURLs, err := getUserImageURLs(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		return nil, err
	}
	if err := commitUserTransaction(ctx, tx); err != nil {
		return nil, err
	}
	return imageURLs, nil
}

func getUserImageURLs(ctx context.Context, tx pgx.Tx, userID uint) ([]string, error) {
	rows, err := tx.Query(ctx, `SELECT image_url FROM user_images WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imageURLs := make([]string, 0)
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, err
		}
		imageURLs = append(imageURLs, imageURL)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return imageURLs, nil
}

func (r *UsersRepository) DeleteUserImage(ctx context.Context, userID, imageID uint) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if err := lockUserRow(ctx, tx, userID); err != nil {
		return "", err
	}

	var imageURL string
	err = tx.QueryRow(ctx, `DELETE FROM user_images
	WHERE id = $1 AND user_id = $2 RETURNING image_url`, imageID, userID).Scan(&imageURL)
	if err != nil {
		return "", err
	}

	imageCount, err := countUserImages(ctx, tx, userID)
	if err != nil {
		return "", err
	}
	if imageCount < 1 {
		return "", ErrLastUserImage
	}

	if err := commitUserTransaction(ctx, tx); err != nil {
		return "", err
	}
	return imageURL, nil
}

func (r *UsersRepository) AddUserImage(ctx context.Context, userID uint, imageURL string) (uint, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	if err := lockUserRow(ctx, tx, userID); err != nil {
		return 0, err
	}

	imageCount, err := countUserImages(ctx, tx, userID)
	if err != nil {
		return 0, err
	}
	if imageCount >= 10 {
		return 0, ErrUserImageLimitExceeded
	}

	var imageID uint
	err = tx.QueryRow(ctx, `INSERT INTO user_images (user_id, image_url)
	VALUES ($1, $2) RETURNING id`, userID, imageURL).Scan(&imageID)
	if err != nil {
		return 0, err
	}

	if err := commitUserTransaction(ctx, tx); err != nil {
		return 0, err
	}
	return imageID, nil
}

func lockUserRow(ctx context.Context, tx pgx.Tx, userID uint) error {
	var id uint
	return tx.QueryRow(ctx, `SELECT id FROM users WHERE id = $1 AND role_id IN (1, 2) FOR UPDATE`, userID).Scan(&id)
}

func countUserImages(ctx context.Context, tx pgx.Tx, userID uint) (int, error) {
	var count int
	err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM user_images WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func commitUserTransaction(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return err
		}
		return errors.Join(ErrUserCommitUncertain, err)
	}
	return nil
}

func (r *UsersRepository) UpdateMyProfile(ctx context.Context, userID uint, user *models.UserProfileUpdateRequest) error {
	return r.UpdateUser(ctx, userID, user, "")
}

func (r *UsersRepository) UpdateUser(ctx context.Context, userID uint, user *models.UserProfileUpdateRequest, passwordHash string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var roleID uint

	err = tx.QueryRow(ctx, `UPDATE users
	SET first_name = $1, gender = $2, bio = $3,
	password_hash = CASE WHEN $4 = '' THEN password_hash ELSE $4 END
	WHERE id = $5 AND role_id IN (1, 2) RETURNING role_id`,
		user.FirstName, user.Gender, user.Bio, passwordHash, userID).Scan(&roleID)

	if err != nil {
		return err
	}

	if roleID == 1 && user.AthleteProfile != nil {
		err = updateAthleteProfile(ctx, tx, userID, user.AthleteProfile)
		if err != nil {
			return err
		}
	} else if roleID == 2 && user.TrainerProfile != nil {
		err = updateTrainerProfile(ctx, tx, userID, user.TrainerProfile)
		if err != nil {
			return err
		}
	}

	if err := updateUserPreferredWorkoutTimes(ctx, tx, userID, user.PreferredWorkoutTimes); err != nil {
		return err
	}
	if err := updateUserGyms(ctx, tx, userID, user.GymIDs); err != nil {
		return err
	}
	return commitUserTransaction(ctx, tx)
}

func updateAthleteProfile(ctx context.Context, tx pgx.Tx, userID uint, profile *models.AthleteProfile) error {
	result, err := tx.Exec(ctx, `
		UPDATE athlete_profiles
		SET date_of_birth = $1::date, height_cm = $2, weight_kg = $3, experience_years = $4,
			min_age = $5, max_age = $6, preferred_gender = $7, min_experience = $8, max_experience = $9
		WHERE user_id = $10`,
		profile.DateOfBirth, profile.Height, profile.Weight, profile.ExperienceYears,
		profile.MinAge, profile.MaxAge, profile.PreferredGender, profile.MinExperience, profile.MaxExperience, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func updateTrainerProfile(ctx context.Context, tx pgx.Tx, userID uint, profile *models.TrainerProfile) error {
	result, err := tx.Exec(ctx, `
		UPDATE trainer_profiles
		SET work_experience_years = $1, education = $2
		WHERE user_id = $3`,
		profile.WorkExperienceYears, profile.Education, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	_, err = tx.Exec(ctx, `DELETE FROM trainer_specializations WHERE trainer_user_id = $1`, userID)
	if err != nil {
		return err
	}

	return addTrainerSpecializations(ctx, tx, userID, profile.SpecializationIDs)
}

func updateUserGyms(ctx context.Context, tx pgx.Tx, userID uint, gymIDs []uint) error {
	if _, err := tx.Exec(ctx, `DELETE FROM user_gyms WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return addUserGyms(ctx, tx, userID, gymIDs)
}

func updateUserPreferredWorkoutTimes(ctx context.Context, tx pgx.Tx, userID uint, preferredTimes []models.PreferredWorkoutTime) error {
	if _, err := tx.Exec(ctx, `DELETE FROM user_preferred_workout_times WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return addUserPreferredWorkoutTimes(ctx, tx, userID, preferredTimes)
}

func (r *UsersRepository) GetUserByID(ctx context.Context, userID uint) (models.UserDetails, error) {
	var user models.UserDetails

	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.email, u.role_id, r.name, u.first_name, u.gender, COALESCE(u.bio, '')
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &user.RoleID, &user.RoleName, &user.FirstName, &user.Gender, &user.Bio)
	if err != nil {
		return models.UserDetails{}, err
	}

	switch user.RoleID {
	case 1:
		profile, err := r.GetAthleteProfile(ctx, userID)
		if err != nil {
			return models.UserDetails{}, err
		}
		user.AthleteProfile = &profile
	case 2:
		profile, err := r.GetTrainerProfile(ctx, userID)
		if err != nil {
			return models.UserDetails{}, err
		}
		user.TrainerProfile = &profile
	}

	user.Images, err = r.GetUserImages(ctx, userID)
	if err != nil {
		return models.UserDetails{}, err
	}
	user.Gyms, err = r.GetUserGyms(ctx, userID)
	if err != nil {
		return models.UserDetails{}, err
	}
	user.PreferredWorkoutTimes, err = r.GetPreferredWorkoutTimes(ctx, userID)
	if err != nil {
		return models.UserDetails{}, err
	}

	return user, nil
}

func (r *UsersRepository) GetAthleteProfile(ctx context.Context, userID uint) (models.AthleteProfile, error) {
	var profile models.AthleteProfile

	err := r.db.QueryRow(ctx, `
		SELECT TO_CHAR(date_of_birth, 'YYYY-MM-DD'), height_cm, weight_kg, experience_years,
			min_age, max_age, preferred_gender, min_experience, max_experience
		FROM athlete_profiles
		WHERE user_id = $1`,
		userID,
	).Scan(&profile.DateOfBirth, &profile.Height, &profile.Weight, &profile.ExperienceYears,
		&profile.MinAge, &profile.MaxAge, &profile.PreferredGender, &profile.MinExperience, &profile.MaxExperience)
	if err != nil {
		return models.AthleteProfile{}, err
	}
	return profile, nil
}

func (r *UsersRepository) GetTrainerProfile(ctx context.Context, userID uint) (models.TrainerProfileDetailed, error) {
	var profile models.TrainerProfileDetailed

	err := r.db.QueryRow(ctx, `
		SELECT work_experience_years, education
		FROM trainer_profiles
		WHERE user_id = $1`,
		userID,
	).Scan(&profile.WorkExperienceYears, &profile.Education)
	if err != nil {
		return models.TrainerProfileDetailed{}, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.name
		FROM trainer_specializations ts
		JOIN specializations s ON s.id = ts.specialization_id
		WHERE ts.trainer_user_id = $1
		ORDER BY s.id`,
		userID,
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

func (r *UsersRepository) GetUserImages(ctx context.Context, userID uint) ([]models.UserImage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, image_url
		FROM user_images
		WHERE user_id = $1
		ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]models.UserImage, 0)
	for rows.Next() {
		var image models.UserImage
		if err := rows.Scan(&image.ID, &image.ImageURL); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
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
