package services

import (
	"context"
	"errors"
	"fitmate/backend/models"
	"fitmate/backend/repositories"
	"fitmate/backend/storage"
	"golang.org/x/crypto/bcrypt"
	"math"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

type AuthService struct {
	authRepository *repositories.AuthRepository
	imageStorage   *storage.ImageStorage
}

func NewAuthService(authRepository *repositories.AuthRepository, imageStorage *storage.ImageStorage) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		imageStorage:   imageStorage,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	registerRequest models.RegisterRequest,
	images []multipart.FileHeader,
) (uint, error) {

	registerRequest.FirstName = strings.TrimSpace(registerRequest.FirstName)
	registerRequest.Email = strings.TrimSpace(registerRequest.Email)

	if registerRequest.TrainerProfile != nil {
		registerRequest.TrainerProfile.Education = strings.TrimSpace(registerRequest.TrainerProfile.Education)
	}

	if len(images) < 1 || len(images) > 10 {
		return 0, ErrImageCountInvalid
	}

	if err := s.validateRegisterRequest(ctx, registerRequest); err != nil {
		return 0, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, ErrPasswordHashFailed
	}

	imageURLs, err := s.imageStorage.SaveUserImages(images)
	if err != nil {
		return 0, err
	}

	userID, err := s.authRepository.CreateUser(ctx, &registerRequest, string(passwordHash), imageURLs)

	if err != nil {
		if !errors.Is(err, repositories.ErrUserCommitUncertain) {
			s.imageStorage.DeleteImages(imageURLs)
		}

		return 0, ErrFailedToCreateUser
	}

	return userID, nil

}

func (s *AuthService) validateRegisterRequest(ctx context.Context, registerRequest models.RegisterRequest) error {
	if err := validateRole(registerRequest.RoleID); err != nil {
		return err
	}

	if err := s.validateEmailForRegister(ctx, registerRequest.Email); err != nil {
		return err
	}
	if err := validatePassword(registerRequest.Password); err != nil {
		return err
	}
	if err := validateFirstName(registerRequest.FirstName); err != nil {
		return err
	}
	if err := validateGender(registerRequest.Gender); err != nil {
		return err
	}
	if err := validateBio(registerRequest.Bio); err != nil {
		return err
	}
	if err := s.validateGyms(ctx, registerRequest.GymIDs); err != nil {
		return err
	}
	if err := validatePreferredWorkoutTimes(registerRequest.PreferredWorkoutTimes); err != nil {
		return err
	}

	if registerRequest.RoleID == 1 {
		if err := validateAthleteProfile(registerRequest.AthleteProfile, registerRequest.TrainerProfile); err != nil {
			return err
		}
	} else if registerRequest.RoleID == 2 {
		if err := s.validateTrainerProfile(ctx, registerRequest.TrainerProfile, registerRequest.AthleteProfile); err != nil {
			return err
		}
	}

	return nil
}

func validateRole(roleID uint) error {
	if roleID != 1 && roleID != 2 {
		return ErrInvalidRole
	}
	return nil
}
func (s *AuthService) validateEmailForRegister(ctx context.Context, email string) error {
	emailLength := utf8.RuneCountInString(email)
	if emailLength == 0 {
		return ErrEmailEmpty
	}
	if emailLength > 100 {
		return ErrEmailTooLong
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return ErrEmailInvalid
	}

	exists, err := s.authRepository.EmailExists(ctx, email)

	if err != nil {
		return ErrEmailFailedCheck
	}
	if exists {
		return ErrEmailExists
	}
	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordNotLongEnough
	}
	return nil
}

func validateFirstName(firstName string) error {
	if firstName == "" {
		return ErrFirstNameEmpty
	}
	if utf8.RuneCountInString(firstName) > 50 {
		return ErrFirstNameTooLong
	}
	return nil
}

func validateGender(gender string) error {
	if utf8.RuneCountInString(gender) == 0 {
		return ErrGenderEmpty
	}

	if gender != "vyras" && gender != "moteris" && gender != "kita" {
		return ErrGenderInvalid
	}
	return nil
}

func validateBio(bio string) error {
	if utf8.RuneCountInString(bio) > 1000 {
		return ErrBioTooLong
	}
	return nil
}

func checkForDuplicateUintSlice(slice []uint) bool {
	seen := make(map[uint]bool)
	for _, value := range slice {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}

func (s *AuthService) validateGyms(ctx context.Context, gymIDs []uint) error {
	checkDublicates := checkForDuplicateUintSlice(gymIDs)
	if checkDublicates {
		return ErrGymIDsInvalid
	}

	for _, gymID := range gymIDs {
		exists, err := s.authRepository.GymExists(ctx, gymID)

		if err != nil {
			return ErrGymIDsFailedCheck
		}
		if !exists {
			return ErrGymIDsInvalid
		}
	}
	return nil
}

func validatePreferredWorkoutTimes(workoutTimes []models.PreferredWorkoutTime) error {
	parsedTimes := make([]models.ParsedPreferredWorkoutTime, 0, len(workoutTimes))

	for _, workoutTime := range workoutTimes {
		parsed, err := parsePreferredWorkoutTime(workoutTime)
		if err != nil {
			return err
		}

		for _, previous := range parsedTimes {
			if previous.DayOfWeek == parsed.DayOfWeek &&
				parsed.StartTime.Before(previous.EndTime) &&
				previous.StartTime.Before(parsed.EndTime) {
				return ErrWorkoutTimesOverlap
			}
		}

		parsedTimes = append(parsedTimes, parsed)
	}
	return nil
}

func parsePreferredWorkoutTime(input models.PreferredWorkoutTime) (models.ParsedPreferredWorkoutTime, error) {
	if input.DayOfWeek < 1 || input.DayOfWeek > 7 {
		return models.ParsedPreferredWorkoutTime{}, ErrWeekdayInvalid
	}

	start, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return models.ParsedPreferredWorkoutTime{}, ErrStartTimeInvalid
	}

	end, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return models.ParsedPreferredWorkoutTime{}, ErrEndTimeInvalid
	}
	if !end.After(start) {
		return models.ParsedPreferredWorkoutTime{}, ErrStartTimeAfterEndTime
	}

	return models.ParsedPreferredWorkoutTime{
		DayOfWeek: input.DayOfWeek,
		StartTime: start,
		EndTime:   end,
	}, nil
}

func validateAthleteProfile(athleteProfile *models.AthleteProfileRequest, trainerProfile *models.TrainerProfileRequest) error {
	if athleteProfile == nil {
		return ErrAthleteProfileEmpty
	}
	if trainerProfile != nil {
		return ErrTrainerProfileNotAllowed
	}

	if err := validateAthleteProfileBirthDate(athleteProfile.DateOfBirth); err != nil {
		return err
	}
	if err := validateAthleteHeight(athleteProfile.Height); err != nil {
		return err
	}
	if err := validateAthleteWeight(athleteProfile.Weight); err != nil {
		return err
	}
	if err := validateExperienceYears(athleteProfile.ExperienceYears); err != nil {
		return err
	}

	return nil
}

func validateExperienceYears(years uint) error {
	if years > math.MaxInt32 {
		return ErrExperienceYearsInvalid
	}
	return nil
}

func validateAthleteHeight(height float64) error {
	if height < 0.01 || height > 999.99 {
		return ErrHeightInvalid
	}
	return nil
}

func validateAthleteWeight(weight float64) error {
	if weight < 0.01 || weight > 999.99 {
		return ErrWeightInvalid
	}
	return nil
}

func validateAthleteProfileBirthDate(birthDate string) error {
	date, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return ErrBirthDateInvalid
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if date.After(today) {
		return ErrBirthDateInvalid
	}

	return nil
}

func (s *AuthService) validateTrainerProfile(ctx context.Context, trainerProfile *models.TrainerProfileRequest, athleteProfile *models.AthleteProfileRequest) error {
	if trainerProfile == nil {
		return ErrTrainerProfileEmpty
	}
	if athleteProfile != nil {
		return ErrAthleteProfileNotAllowed
	}

	if err := validateTrainerEducation(trainerProfile.Education); err != nil {
		return err
	}
	if err := validateExperienceYears(trainerProfile.WorkExperienceYears); err != nil {
		return err
	}
	if err := s.validateTrainerSpecializations(ctx, trainerProfile.SpecializationIDs); err != nil {
		return err
	}
	return nil
}

func validateTrainerEducation(education string) error {
	if education == "" {
		return ErrEducationEmpty
	}
	if utf8.RuneCountInString(education) > 255 {
		return ErrEducationTooLong
	}
	return nil
}

func (s *AuthService) validateTrainerSpecializations(ctx context.Context, specializationIDs []uint) error {
	checkDublicates := checkForDuplicateUintSlice(specializationIDs)
	if checkDublicates {
		return ErrSpecializationIDsInvalid
	}

	for _, specializationID := range specializationIDs {
		exists, err := s.authRepository.SpecializationExists(ctx, specializationID)

		if err != nil {
			return ErrSpecializationIDsFailedCheck
		}
		if !exists {
			return ErrSpecializationIDsInvalid
		}
	}

	return nil
}
