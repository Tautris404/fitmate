package services

import (
	"context"
	"errors"
	"math"
	"mime/multipart"
	"os"
	"strings"

	"fitmate/backend/models"
	"fitmate/backend/repositories"
	"fitmate/backend/storage"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
	authRepository  *repositories.AuthRepository
	imageStorage    *storage.ImageStorage
}

func NewUsersService(usersRepository *repositories.UsersRepository, authRepository *repositories.AuthRepository, imageStorage *storage.ImageStorage) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		authRepository:  authRepository,
		imageStorage:    imageStorage,
	}
}

func (s *UsersService) GetUsers(ctx context.Context, page, limit int) (models.UserListResponse, error) {
	if err := validatePagination(page, limit); err != nil {
		return models.UserListResponse{}, err
	}

	offset := (page - 1) * limit
	users, total, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return models.UserListResponse{}, ErrFailedToGetUsers
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	return models.UserListResponse{
		Users:      users,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func validatePagination(page, limit int) error {
	if page < 1 {
		return ErrInvalidPage
	}
	if limit < 1 || limit > 100 {
		return ErrInvalidLimit
	}
	if page-1 > math.MaxInt/limit {
		return ErrInvalidPage
	}
	return nil
}

func (s *UsersService) DeleteUser(ctx context.Context, userID uint) error {
	if err := validateUserID(userID); err != nil {
		return ErrProfileNotFound
	}

	imageURLs, err := s.usersRepository.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProfileNotFound
		}
		return ErrFailedToDeleteProfile
	}
	if err := s.imageStorage.DeleteImages(imageURLs); err != nil {
		return ErrFailedToDeleteProfile
	}
	return nil
}

func (s *UsersService) AddUserImage(ctx context.Context, userID uint, image *multipart.FileHeader) (uint, error) {
	if err := validateUserID(userID); err != nil {
		return 0, ErrProfileNotFound
	}
	if image == nil {
		return 0, ErrSingleImageRequired
	}

	imageURL, err := s.imageStorage.SaveUserImage(image)
	if err != nil {
		return 0, err
	}

	imageID, err := s.usersRepository.AddUserImage(ctx, userID, imageURL)
	if err != nil {
		if !errors.Is(err, repositories.ErrUserCommitUncertain) {
			s.imageStorage.DeleteImage(imageURL)
		}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return 0, ErrProfileNotFound
		case errors.Is(err, repositories.ErrUserImageLimitExceeded):
			return 0, ErrImageLimitExceeded
		default:
			return 0, ErrFailedToAddImage
		}
	}
	return imageID, nil
}

func (s *UsersService) DeleteUserImage(ctx context.Context, userID, imageID uint) error {
	if err := validateUserID(userID); err != nil {
		return ErrProfileNotFound
	}
	if err := validateImageID(imageID); err != nil {
		return err
	}

	imageURL, err := s.usersRepository.DeleteUserImage(ctx, userID, imageID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrImageNotFound
		case errors.Is(err, repositories.ErrLastUserImage):
			return ErrLastImage
		default:
			return ErrFailedToDeleteImage
		}
	}
	if err := s.imageStorage.DeleteImage(imageURL); err != nil && !errors.Is(err, os.ErrNotExist) {
		return ErrFailedToDeleteImage
	}
	return nil
}

func validateImageID(imageID uint) error {
	if imageID == 0 || imageID > math.MaxInt32 {
		return ErrInvalidImageID
	}
	return nil
}

func (s *UsersService) GetMyProfile(ctx context.Context, userID uint) (models.UserProfile, error) {
	if err := validateUserID(userID); err != nil {
		return models.UserProfile{}, ErrProfileNotFound
	}

	profile, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserProfile{}, ErrProfileNotFound
		}
		return models.UserProfile{}, ErrFailedToGetProfile
	}

	return profile.UserProfile, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, userID uint) (models.UserDetails, error) {
	if err := validateUserID(userID); err != nil {
		return models.UserDetails{}, ErrProfileNotFound
	}

	user, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserDetails{}, ErrProfileNotFound
		}
		return models.UserDetails{}, ErrFailedToGetProfile
	}
	return user, nil
}

func (s *UsersService) UpdateMyProfile(ctx context.Context, userID uint, user *models.UserProfileUpdateRequest) error {
	if err := validateUserID(userID); err != nil {
		return err
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	if user.TrainerProfile != nil {
		user.TrainerProfile.Education = strings.TrimSpace(user.TrainerProfile.Education)
	}

	profile, err := s.GetMyProfile(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.validateProfileUpdateRequest(ctx, profile.RoleID, user); err != nil {
		return err
	}

	if err := s.usersRepository.UpdateMyProfile(ctx, userID, user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProfileNotFound
		}
		return ErrFailedToUpdateProfile
	}
	return nil
}

func (s *UsersService) UpdateUser(ctx context.Context, userID uint, user *models.UserUpdateRequest) error {
	profile, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	if user.TrainerProfile != nil {
		user.TrainerProfile.Education = strings.TrimSpace(user.TrainerProfile.Education)
	}
	if err := s.validateProfileUpdateRequest(ctx, profile.RoleID, &user.UserProfileUpdateRequest); err != nil {
		return err
	}

	passwordHash := ""
	if user.Password != "" {
		if err := validatePassword(user.Password); err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				return ErrPasswordTooLong
			}
			return ErrPasswordHashFailed
		}
		passwordHash = string(hash)
	}

	if err := s.usersRepository.UpdateUser(ctx, userID, &user.UserProfileUpdateRequest, passwordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProfileNotFound
		}
		return ErrFailedToUpdateProfile
	}
	return nil
}

func (s *UsersService) validateProfileUpdateRequest(ctx context.Context, roleID uint, user *models.UserProfileUpdateRequest) error {
	if err := validateRole(roleID); err != nil {
		return err
	}
	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}
	if err := validateGender(user.Gender); err != nil {
		return err
	}
	if err := validateBio(user.Bio); err != nil {
		return err
	}
	if err := validateGyms(ctx, s.authRepository, user.GymIDs); err != nil {
		return err
	}
	if err := validatePreferredWorkoutTimes(user.PreferredWorkoutTimes); err != nil {
		return err
	}

	if roleID == 1 {
		return validateAthleteProfile(user.AthleteProfile, user.TrainerProfile)
	}
	return validateTrainerProfile(ctx, s.authRepository, user.TrainerProfile, user.AthleteProfile)
}

func validateUserID(userID uint) error {
	if userID == 0 || userID > math.MaxInt32 {
		return ErrInvalidUserID
	}
	return nil
}
